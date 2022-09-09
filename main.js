//@ts-check

const CLI = require('./cli.js');
const AWSClient = require('./aws');

CLI.Setup(process.argv);

/**
 * debugInfo prints the params and fileParams to the console
 * @param {Record<string, any>} options CLI options
 */
function debugInfo(options) {
	console.log("=================")
	console.log(options);
	console.log(AWSClient.prepareParams(options.params, options.fp));
	console.log("=================")
}

async function main() {
	CLI.Setup(process.argv);
	const options = CLI.Options();

	if (options.debug) {
		debugInfo(options);
		process.exit(0);
	}

	// @ts-ignore
	const client = new AWSClient(options);
	const res = await client.runAny(options.apiName, options.params, options.fp);
	console.log(res);
}

main();