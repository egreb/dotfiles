#!/usr/bin/env python3
"""Reconstructed integration suite: local Git, fake Docker CLI/tools, isolated tmux."""
import fcntl, json, os, pathlib, pty, re, select, socket as sockets, struct, subprocess, tempfile, termios, time
P=pathlib.Path
binary=str(P(os.environ.get('WORKFLOW','dist/sbx')).resolve())
with tempfile.TemporaryDirectory(prefix='sbx-recovery-test-') as tmp:
 root=P(tmp);fake=root/'bin';fake.mkdir();cfg=root/'config';cfg.mkdir();work=root/'workspaces';work.mkdir();state=root/'state';state.mkdir()
 socket='sbx-recovery-'+str(os.getpid())
 env=dict(os.environ,HOME=str(root),SBX_CONFIG_DIR=str(cfg),SBX_CONFIG_FILE=str(cfg/'config.yaml'),SBX_FAKE_STATE=str(state),PATH=str(fake)+os.pathsep+os.environ['PATH'],TERM='xterm-256color')
 env.pop('TMUX',None);env.pop('TMUX_PANE',None)
 env.pop('SBX_TMUX_PANE',None);env.pop('SBX_TMUX_CLIENT',None)
 def run(args,ok=True,**kw):
  r=subprocess.run(args,env=env,text=True,capture_output=True,**kw)
  if ok and r.returncode:raise AssertionError(f'{args}: {r.stdout}\n{r.stderr}')
  return r
 def tmux(*args,ok=True):return run(['tmux','-L',socket,*args],ok=ok)
 def sbx(*args,ok=True):return run([binary,*args],ok=ok)
 def executable(name,contents):p=fake/name;p.write_text(contents);p.chmod(0o755)
 executable('native-sbx', '''#!/usr/bin/env python3
import os,sys,json,time,pathlib
p=pathlib.Path(os.environ['SBX_FAKE_STATE']);args=sys.argv[1:]
with (p/'calls.jsonl').open('a') as f:f.write(json.dumps(args)+'\\n')
cmd=args[0]
if cmd=='ls':
 if args!=['ls','--quiet']:sys.exit('unsupported ls arguments: '+repr(args))
 for x in p.glob('*.sandbox'):print(x.stem)
elif cmd=='create':
 name=args[args.index('--name')+1];(p/(name+'.sandbox')).touch()
 if name=='agentfail':sys.exit(7)
elif cmd=='rm':(p/(args[-1]+'.sandbox')).unlink(missing_ok=True)
elif cmd=='run':time.sleep(120)
elif cmd=='exec':pass
elif cmd=='version':print('fake native')
else:sys.exit(2)
''')
 for tool in ['nvim','lazygit','hunk']:executable(tool,'#!/bin/sh\nprintf "tool: %s\\n" "'+tool+'"\n')
 source=root/'repository';source.mkdir();run(['git','init','-b','main',str(source)])
 (source/'README.md').write_text('local integration fixture\n');run(['git','-C',str(source),'add','.']);run(['git','-C',str(source),'-c','user.name=Test','-c','user.email=test@example.invalid','commit','-m','fixture'])
 (cfg/'projects.yaml').write_text('projects:\n  - name: backend\n    repo: '+json.dumps(str(source))+'\n  - name: frontend\n    repo: '+json.dumps(str(source))+'\n')
 (cfg/'config.yaml').write_text('workspace_root: '+json.dumps(str(work))+'\nprojects_file: '+json.dumps(str(cfg/'projects.yaml'))+'\ngit_cache_root: '+json.dumps(str(root/'cache'))+'\nhost_shell: /bin/sh\nagent: codex\nclaude_theme: ""\ntmux_socket: '+socket+'\ntmux_config: '+json.dumps(str(cfg/'tmux.conf'))+'\ndocker_sbx_bin: '+json.dumps(str(fake/'native-sbx'))+'\n')
 (cfg/'tmux.conf').write_text('set -g prefix C-s\nset -g status off\nset -g base-index 1\nbind r source-file '+str(cfg/'tmux.conf')+'\nbind N display-message work-notes\n')
 master=slave=None;client=None
 try:
  sbx('doctor');sbx('new','--all','--no-open','alpha')
  # Exercise the shared navigation config when testing inside dotfiles.
  global_config=P(__file__).resolve().parents[2]/'.tmux.conf'
  if global_config.is_file():
   original=(cfg/'tmux.conf').read_text()
   (cfg/'tmux.conf').write_text('source-file '+json.dumps(str(global_config))+'\n'+original)
   sbx('tmux-reload')
  assert (work/'alpha'/'.sbx-managed').is_file()
  assert run(['git','-C',str(work/'alpha'/'backend'),'remote','get-url','origin']).stdout.strip()==str(source)
  assert '__popup tui' in tmux('list-keys','-T','prefix','o').stdout
  assert '__popup new' in tmux('list-keys','-T','prefix','i').stdout
  assert 'work-notes' in tmux('list-keys','-T','prefix','N').stdout
  print('PASS creation, local cache clones, metadata, overlay preserves notes binding',flush=True)
  master,slave=pty.openpty()
  fcntl.ioctl(slave,termios.TIOCSWINSZ,struct.pack('HHHH',40,140,0,0))
  client=subprocess.Popen(['tmux','-L',socket,'attach-session','-t','alpha'],env=env,stdin=slave,stdout=slave,stderr=slave)
  for _ in range(30):
   if tmux('list-clients').stdout.strip():break
   time.sleep(.05)
  def terminal_until(needle,timeout=5):
   output=b'';deadline=time.monotonic()+timeout
   while time.monotonic()<deadline:
    if select.select([master],[],[],.1)[0]:
     output+=os.read(master,65536)
     if needle in output:return output
   raise AssertionError(f'Expected {needle!r} from tmux client; last output: {output[-1500:]!r}')
  def drain_terminal():
   output=b''
   while select.select([master],[],[],.05)[0]:output+=os.read(master,65536)
   return output
  drain_terminal()
  # Exercise the real keybinding: popups do not automatically inherit TMUX_PANE.
  os.write(master,b'\x13f')
  terminal_until(b'Open: alpha / agent')
  # Select backend, then reopen the same binding from that project and return to agent.
  for keys,target,next_identity in [(b'\x1b[B\r','alpha--backend',b'Open: alpha / backend'),(b'\x1b[A\r','alpha',None)]:
   drain_terminal();os.write(master,keys)
   interaction=b''
   for _ in range(80):
    interaction+=drain_terminal()
    if tmux('list-clients','-F','#{session_name}').stdout.strip()==target:break
    time.sleep(.05)
   actual=tmux('list-clients','-F','#{session_name}').stdout.strip()
   if actual!=target:
    text=re.sub(rb'\x1b\[[0-?]*[ -/]*[@-~]|\x1b[()][A-Za-z0-9]',b'',interaction).decode(errors='replace')
    text=re.sub(r'[│─┌┐└┘\s]+',' ',text)
    raise AssertionError((actual,target,text[-2000:]))
   if next_identity:
    drain_terminal();os.write(master,b'\x13f');terminal_until(next_identity)
  print('PASS actual prefix+f popup resolves agent/project context and switches original client',flush=True)
  pane=tmux('display-message','-p','-t','alpha:','#{pane_id}').stdout.strip();sp=tmux('display-message','-p','#{socket_path}').stdout.strip()
  env['TMUX']=sp+',1,0';env['TMUX_PANE']=pane
  sbx('project','backend');sbx('project','backend')
  windows=tmux('list-windows','-t','alpha--backend','-F','#{window_index}:#{window_name}:#{window_panes}').stdout.strip()
  assert windows=='1:neovim:1\n2:shell:1',windows
  assert tmux('show-options','-qv','-t','alpha','@sbx-last-session').stdout.strip()=='alpha--backend'
  print('PASS project sessions, Neovim and shell windows, idempotent reuse',flush=True)
  # Split from the original agent pane without changing its last-session choice.
  tmux('switch-client','-t','alpha');sbx('split','frontend')
  panes=tmux('list-panes','-t','alpha:','-F','#{pane_id}').stdout.splitlines();side=next(x for x in panes if x!=pane)
  assert tmux('show-options','-p','-qv','-t',side,'@sbx-sidecar').stdout.strip()=='1'
  assert tmux('show-options','-qv','-t','alpha','@sbx-last-session').stdout.strip()=='alpha--backend'
  for _ in range(40):
   r=sbx('__sidecar-key','2',side,ok=False)
   if r.returncode==0:break
   time.sleep(.05)
  assert r.returncode==0,r.stderr
  # Navigation must stay in the outer window when a nested sidecar client is focused.
  for key,wanted in [(b'\x08',pane),(b'\x0c',side),(b'\x08',pane)]:
   drain_terminal();os.write(master,key)
   for _ in range(30):
    drain_terminal()
    active=tmux('display-message','-p','-t','alpha:','#{pane_id}').stdout.strip()
    if active==wanted:break
    time.sleep(.05)
   assert active==wanted,('pane navigation',key,active,wanted)
  print('PASS Ctrl-h/l between outer agent and nested sidecar panes',flush=True)
  tmux('kill-pane','-t',side);tmux('has-session','-t','alpha--frontend')
  print('PASS sidecar attachment, numeric keys, closing preserves target session',flush=True)
  # The same keys should work in ordinary splits and while in vi copy mode.
  regular=tmux('split-window','-h','-P','-F','#{pane_id}','-t',pane,'/bin/sh').stdout.strip()
  for key,wanted in [(b'\x08',pane),(b'\x0c',regular)]:
   drain_terminal();os.write(master,key)
   for _ in range(30):
    drain_terminal()
    active=tmux('display-message','-p','-t','alpha:','#{pane_id}').stdout.strip()
    if active==wanted:break
    time.sleep(.05)
   assert active==wanted,('ordinary navigation',key,active,wanted)
  if global_config.is_file():
   tmux('set-window-option','-t','alpha:','mode-keys','vi')
   tmux('copy-mode','-t',regular);drain_terminal();os.write(master,b'\x08')
   for _ in range(30):
    drain_terminal()
    active=tmux('display-message','-p','-t','alpha:','#{pane_id}').stdout.strip()
    if active==pane:break
    time.sleep(.05)
   assert active==pane,('copy mode navigation',active,pane)
  tmux('kill-pane','-t',regular)
  print('PASS Ctrl-h/l in ordinary splits and vi copy mode',flush=True)

  # Reload twice must preserve bindings and the existing sessions.
  sbx('tmux-reload');sbx('tmux-reload')
  assert 'work-notes' in tmux('list-keys','-T','prefix','N').stdout
  assert '__sidecar-key' in tmux('list-keys','-T','prefix','2').stdout
  tmux('switch-client','-t','alpha--backend');tmux('kill-session','-t','alpha')
  env['TMUX_PANE']=tmux('display-message','-p','-t','alpha--backend:','#{pane_id}').stdout.strip()
  sbx('open','alpha')
  assert tmux('show-options','-qv','-t','alpha','@sbx-role').stdout.strip()=='agent'
  assert tmux('show-options','-qv','-t','alpha','@sbx-last-session').stdout.strip()=='alpha--backend'
  # Add only a missing configured project, preserving both existing clones.
  with (cfg/'projects.yaml').open('a') as f:f.write('  - name: worker\n    repo: '+json.dumps(str(source))+'\n')
  sbx('add','worker');assert (work/'alpha'/'worker'/'.git').is_dir()
  assert sbx('add','worker',ok=False).returncode
  print('PASS repeated reload, missing-agent recovery, adding a project',flush=True)
  # Passive list must not contact Docker.
  calls=(state/'calls.jsonl').read_text();assert sbx('list','--names').stdout.strip()=='alpha';assert calls==(state/'calls.jsonl').read_text()
  assert sbx('new','--all','--no-open','../outside',ok=False).returncode
  (root/'outside').mkdir();(root/'outside'/'.sbx-managed').write_text('format=1\nname=escape\n');(work/'escape').symlink_to(root/'outside',target_is_directory=True)
  assert sbx('delete','--yes','escape',ok=False).returncode
  (work/'unmarked').mkdir();assert sbx('delete','--yes','unmarked',ok=False).returncode
  assert (root/'outside').is_dir() and (work/'unmarked').is_dir()
  print('PASS passive listing, traversal/symlink/unmarked deletion rejection',flush=True)
  # Failure after native create must roll back resources reserved by this command.
  assert sbx('new','--all','--no-open','agentfail',ok=False).returncode
  assert not (work/'agentfail').exists() and not (state/'agentfail.sandbox').exists()
  sbx('new','--all','--no-open','beta')
  # Delete from the normal tree picker, including its current workspace.
  sbx('new','--all','--no-open','picker-delete')
  time.sleep(1.1);tmux('switch-client','-t','beta');time.sleep(1.1)
  tmux('switch-client','-t','picker-delete')
  drain_terminal();os.write(master,b'\x13o');terminal_until(b'sbx sessions')
  os.write(master,b'D');terminal_until(b'Delete workspace picker-delete?')
  os.write(master,b'n');time.sleep(.2)
  assert (work/'picker-delete').is_dir(), 'cancel deleted workspace'
  os.write(master,b'D');terminal_until(b'Delete workspace picker-delete?')
  os.write(master,b'y');terminal_until(b'Create a new workspace?')
  assert not (work/'picker-delete').exists(), 'picker delete failed'
  os.write(master,b'n');time.sleep(.3)
  active=tmux('list-clients','-F','#{session_name}').stdout.strip()
  assert active=='beta', 'did not return to last used workspace: '+active
  print('PASS tree picker deletion cancellation, current-workspace deletion, and last-used fallback',flush=True)
  # A host dev server that ignores terminal hangup and graceful termination.
  server=root/'server.py'
  server.write_text("import socket,signal,pathlib,time,sys\nsignal.signal(signal.SIGHUP,signal.SIG_IGN)\nsignal.signal(signal.SIGTERM,signal.SIG_IGN)\ns=socket.socket();s.bind(('127.0.0.1',0));s.listen()\npathlib.Path(sys.argv[1]).write_text(str(s.getsockname()[1]))\ntime.sleep(120)\n")
  import shlex
  for workspace in ['alpha','beta']:
   tmux('new-window','-d','-t',workspace,'python3 '+shlex.quote(str(server))+' '+shlex.quote(str(root/(workspace+'.port'))))
  for _ in range(100):
   if all((root/(w+'.port')).exists() for w in ['alpha','beta']):break
   time.sleep(.05)
  alpha_port=int((root/'alpha.port').read_text());beta_port=int((root/'beta.port').read_text())
  tmux('switch-client','-t','alpha')
  drain_terminal();os.write(master,b'\x13D');terminal_until(b'Delete workspace: alpha')
  os.write(master,b'\r');terminal_until(b'Type the workspace name:')
  assert (work/'alpha').is_dir() and (state/'alpha.sandbox').exists()
  os.write(master,b'wrong-name\r');time.sleep(.2)
  assert (work/'alpha').is_dir() and (state/'alpha.sandbox').exists()
  drain_terminal();os.write(master,b'\x13D');terminal_until(b'Delete workspace: alpha')
  os.write(master,b'\r');terminal_until(b'Type the workspace name:')
  os.write(master,b'alpha\r')
  for _ in range(160):
   if not (work/'alpha').exists():break
   time.sleep(.05)
  assert not (work/'alpha').exists(), 'current workspace deletion did not complete'
  print('PASS delete popup from current workspace, wrong-name cancellation, exact-name confirmation',flush=True)
  with sockets.socket() as probe:probe.bind(('127.0.0.1',alpha_port))
  with sockets.create_connection(('127.0.0.1',beta_port),timeout=1):pass
  sbx('delete','--yes','beta')
  with sockets.socket() as probe:probe.bind(('127.0.0.1',beta_port))
  sbx('new','--all','--no-open','beta')
  print('PASS deletion releases host server port and preserves other workspace processes',flush=True)
  assert not (work/'alpha').exists() and (work/'beta').is_dir()
  assert not (state/'alpha.sandbox').exists() and (state/'beta.sandbox').exists()
  assert 'alpha' not in tmux('list-sessions','-F','#{session_name}').stdout
  print('PASS failed creation rollback and exact workspace deletion',flush=True)
  # The last workspace must leave the popup alive for creation, too.
  client.wait(timeout=3)
  client=subprocess.Popen(['tmux','-L',socket,'attach-session','-t','beta'],env=env,stdin=slave,stdout=slave,stderr=slave)
  for _ in range(30):
   if tmux('list-clients').stdout.strip():break
   time.sleep(.05)
  drain_terminal();os.write(master,b'\x13o');terminal_until(b'sbx sessions')
  os.write(master,b'D');terminal_until(b'Delete workspace beta?')
  os.write(master,b'y');terminal_until(b'Create a new workspace?')
  assert not (work/'beta').exists()
  os.write(master,b'y');terminal_until(b'New workspace')
  os.write(master,b'\x1b');time.sleep(.3)
  active=tmux('list-clients','-F','#{session_name}').stdout.strip()
  assert active.startswith('sbx-picker-'), 'last deletion did not preserve a host shell: '+active
  assert tmux('show-options','-qv','-t',active,'@sbx-managed').stdout.strip()!='1'
  print('PASS last-workspace deletion, creation prompt, and usable shell after canceled creation',flush=True)
 finally:
  tmux('kill-server',ok=False)
  if client:
   try:client.wait(timeout=3)
   except subprocess.TimeoutExpired:client.kill();client.wait(timeout=3)
  if master is not None:os.close(master)
  if slave is not None:os.close(slave)
