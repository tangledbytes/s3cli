//@ts-check
const package = require('./package.json');

module.exports = {
	Version: function() {
		return package.version;
	},
	Desc: function() {
		return package.description;
	}
}