-- 'literals' only hints parameter names for literal args (foo(3) -> foo(retries: 3));
-- change to 'all' to show them at every call site
local inlay_hints = {
  parameterNames = { enabled = 'literals', suppressWhenArgumentMatchesName = true },
  parameterTypes = { enabled = true },
  variableTypes = { enabled = true },
  propertyDeclarationTypes = { enabled = true },
  functionLikeReturnTypes = { enabled = true },
  enumMemberValues = { enabled = true },
}

return {
  filetypes = { 'typescript', 'javascript', 'typescriptreact', 'javascriptreact' },
  root_markers = { 'tsconfig.json', 'jsconfig.json', 'package.json', '.git' },
  settings = {
    typescript = { inlayHints = inlay_hints },
    javascript = { inlayHints = inlay_hints },
  },
}
