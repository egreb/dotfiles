return {
  settings = {
    gopls = {
      staticcheck = true,
      analyses = {
        unusedparams = true,
        unusedwrite = true,
        unusedvariable = true,
        shadow = true,
        nilness = true,
        useany = true,
      },
      hints = {
        parameterNames = true,
        assignVariableTypes = true,
        compositeLiteralFields = true,
        constantValues = true,
        rangeVariableTypes = true,
      },
      vulncheck = 'Imports',
      diagnosticsDelay = '300ms',
      usePlaceholders = true,
    },
  },
}
