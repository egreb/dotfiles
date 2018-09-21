#!/usr/bin/env python3
import os
import subprocess
import argparse

FNULL = open(os.devnull, 'w')
deleteAll = False
deleteMerged = False

def main():
	parseArgs()
	branchCmd = subprocess.run(['git', 'branch'], capture_output=True)
	if branchCmd.returncode == 1:
		print('You dont seem to be in a git repository')
		exit()

	exclude = ['*', 'master']
	branches = [branch for branch in branchCmd.stdout.decode().split() if branch not in exclude]
	if len(branches) is 0:
		print('No branches requires attention.')
		exit()

	for branch in branches:
			status = subprocess.run(['git', 'branch', '-d', branch], stderr=FNULL)
			if status.returncode == 1:
				if not deleteAll:
					print('Could not delete branch ' + branch +
								' since its not merged with master, still wanna delete it?')
					response = input('no/yes/all/abort: ')
					handleResponse(response, branch)
				else:
					deleteBranch(branch)
			else:
				print(branch + ' deleted')

	print('all done :)')


def handleResponse(response, branch):
	if response == "abort" or response == 'a':
		print('aborting')
		exit()
	if response == 'all':
		deleteAll = True
	if response == 'yes' or response == 'y':
		deleteBranch(branch)

def deleteBranch(branch):
	status = subprocess.run(['git', 'branch', '-D', branch], stderr=FNULL, stdout=FNULL)
	if status.returncode != 1:
			print(branch + ' is deleted')
	else:
		print('Could not delete ' + branch)

def parseArgs():
	parser = argparse.ArgumentParser(
		description='Cleanup git repository branches')
	parser.add_argument('-a', '--all', default=False, help='Delete all branches without prompt')
	parser.add_argument('-m', '--merged', default=False, help='Delete all merged branches')
	args = parser.parse_args()

	if args.all:
		deleteAll = True
	else:
		deleteAll = False

	if args.merged:
		deleteMerged = True
	else:
		deleteMerged = False
main()
