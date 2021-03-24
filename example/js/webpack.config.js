const path = require('path');
const WebpackObfuscator = require('webpack-obfuscator');

module.exports = (env = {}) => {
    const { ENTRY, SHOULD_OBFUSCATE, FILENAME = 'index.js' } = env;
    const OBFUSCATE = SHOULD_OBFUSCATE === 'true';

    const obfuscatePlugin = OBFUSCATE ?
        [
            new WebpackObfuscator ({
                rotateStringArray: true
            }, [])
        ]
        : []


    return {
        mode: 'production',
        entry: ENTRY.split(','),
        output: {
            path: path.resolve(__dirname + `/../out/js/`),
            filename: `${FILENAME}.js`,
          },
        plugins: obfuscatePlugin,
    }
}