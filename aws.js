//@ts-check
const AWS = require("@aws-sdk/client-s3");
const { NodeHttpHandler } = require('@aws-sdk/node-http-handler')
const http = require('http');
const https = require('https');
const fs = require('fs');

/**
 * credentials returns a credentials object if accessKey and secretKey are provided
 * or else returns undefined
 * @param {{
 * 	anon?: boolean,
 * 	accessKey?: string,
 * 	secretKey?: string,
 * }} options 
 * @returns 
 */
function credentials(options) {
	if (options.anon) return undefined;

	return {
		accessKeyId: options.accessKey,
		secretAccessKey: options.secretKey,
	};
}

class AWSClient {
	/**
	 * Generic AWS S3 client
	 * @param {{
	 * 	secretKey?: string,
	 * 	accessKey?: string,
	 * 	skipSsl?: boolean,
	 * 	endpoint: string,
	 * 	anon?: boolean,
	 * 	tls?: boolean,
	 * }} options 
	 */
	constructor(options) {
		this.options = options;

		const handler = new NodeHttpHandler({
			httpAgent: new http.Agent({ keepAlive: false }),
			httpsAgent: new https.Agent({ keepAlive: false, rejectUnauthorized: !this.options.skipSsl }),
		});

		this.client = new AWS.S3Client({
			region: "us-east-1",
			// @ts-ignore
			credentials: credentials(this.options),
			endpoint: this.options.endpoint,
			requestHandler: handler,
			forcePathStyle: true,
			tls: !!this.options.tls,
		});
	}

	/**
	 * runAny takes an API name and runs that API with the given params
	 * @param {string} name API Name
	 * @param {Record<string, any>} params API Parameters
	 * @param {Record<string, any>} fileParams API File Parameters
	 * @returns 
	 */
	async runAny(name, params, fileParams) {
		const Command = AWS[`${name}Command`];
		if (!Command) throw new Error(`Command ${name} not found`);

		const finalParams = AWSClient.prepareParams(params, fileParams);

		const command = new Command(finalParams);
		return this.client.send(command);
	}

	/**
	 * prepareParams takes a params object and a fileParams object and returns a new params object
	 * with the fileParams added to it
	 * @param {Record<string, any>} params 
	 * @param {Record<string, any>} fileParams 
	 * @returns {Record<string, any>}
	 */
	static prepareParams(params, fileParams) {
		const fp = fileParams || {};
		const newParams = Object.assign({}, params);

		Object.keys(fp).forEach(key => {
			const value = fp[key];
			newParams[key] = fs.readFileSync(value, 'utf8');
		});

		return newParams;
	}
}

module.exports = AWSClient;