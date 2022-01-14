const path = require('path')
const microargs = require('microargs')
const { get, difference, isEmpty, padEnd, forEach, capitalize, omit, isString } = require('lodash')

class CLIError extends Error { }

function optionToString (optionName) {
  return optionName.length === 1 ? `-${optionName}` : `--${optionName}`
}

function optionsToString (optionsKeys) {
  return optionsKeys.map(optionToString).join(' ')
}

function printHelp (scriptName, annotations, logger) {
  if (isEmpty(annotations)) {
    logger.log('Documentation not found')
    return null
  }

  const { description, params, options } = annotations
  const extra = omit(annotations, ['description', 'params', 'options'])
  const usageOptions = isEmpty(options) ? '' : '[options]'
  const usageParams = isEmpty(params) ? '' : `[${params.join(' ')}]`

  logger.log(`Usage: ${path.basename(scriptName)} ${usageOptions} ${usageParams}\n`)

  if (description) {
    logger.log(`${description}\n`)
  }

  if (!isEmpty(options)) {
    logger.log('Options:\n')
    forEach(options, (value, key) => {
      logger.log(`  ${padEnd(optionToString(key), 12)}${value}`)
    })
  }

  forEach(extra, (value, key) => {
    logger.log(`\n${capitalize(key)}:\n`)
    logger.log(`${value}\n`)
  })
}

module.exports = (argv, annotations = {}, help, logger = console) => {
  help = help || printHelp

  return (callback) => {
    const { params, options } = microargs(argv.slice(2))
    const scriptName = path.basename(argv[1])

    if (isString(annotations)) {
      annotations = {
        description: annotations
      }
    }

    if (options.help) {
      return help(scriptName, annotations, logger)
    }

    const annotatedOptionsKeys = (get(annotations, 'options') && Object.keys(annotations.options)) || []
    const optionsKeys = Object.keys(options)
    const illegalOptionsKeys = difference(optionsKeys, annotatedOptionsKeys)

    if (annotatedOptionsKeys.length && illegalOptionsKeys.length) {
      const msg = `Illegal option: ${optionsToString(illegalOptionsKeys)}\n` +
        `Available options: ${optionsToString(annotatedOptionsKeys)}\n` +
        `Type "${scriptName} --help" for more information`
      throw new CLIError(msg)
    }

    return callback(options, ...params)
  }
}

module.exports.CLIError = CLIError
