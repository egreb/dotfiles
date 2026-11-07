-- 'literals' only hints parameter names for literal args (foo(3) -> foo(retries: 3));
-- change to 'all' to show them at every call site
local inlay_hints = {
  parameterNames = { enabled = 'literals', suppressWhenArgumentMatchesName = true },
  parameterTypes = { enabled = false },
  variableTypes = { enabled = false },
  propertyDeclarationTypes = { enabled = false },
  functionLikeReturnTypes = { enabled = false },
  enumMemberValues = { enabled = false },
}

return {
  filetypes = { 'typescript', 'javascript', 'typescriptreact', 'javascriptreact' },
  root_markers = { 'tsconfig.json', 'jsconfig.json', 'package.json', '.git' },
  settings = {
    typescript = { inlayHints = inlay_hints },
    javascript = { inlayHints = inlay_hints },
  },
}
