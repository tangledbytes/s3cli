//@ts-check
const { program } = require('commander');

/**
 * Setup takes in the argument array and and uses that setup the commander program
 * @param {Array<string>} args arguments
 */
function Setup(args) {
	program
		.requiredOption('--endpoint <endpoint>', 'endpoint')
		.requiredOption('--api <apiName>', 'api name')
		.option('--bucket <bucket>', 'bucket name')
		.option('--accessKey <accessKey>', 'access key')
		.option('--secretKey <secretKey>', 'secret key')
		.option('--params <params>', 'params')
		.option('--fp <fp>', 'fp is file parameter')
		.option('--tls <tls>', 'tls', false)
		.option('--anon', 'anonymous', false)
		.option('--skip-ssl', 'skip ssl', false)
		.option('--debug', 'debug')
		.parse(args);
}

/**
 * Options return the object of options
 * @returns {Record<string, any>}
 */
function Options() {
	const opts = program.opts();

	opts.params = JSON.parse(opts.params || '{}');
	opts.fp = JSON.parse(opts.fp || '{}');

	return opts;
}

module.exports = {
	Setup,
	Options,
}