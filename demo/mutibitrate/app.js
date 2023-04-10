const crypto = require('crypto');
const http = require('http');

host = 'api.agora.io';

appid = '552cbXXXXXXXXXXXxef8f13';
username = '10eXXXXXXXXXXXXX2627';
secret = 'e947XXXXXXXXXXXXXXd';

function signData(data) {
    const hmac = crypto.createHmac('sha256', secret);
    hmac.update(data);
    return hmac.digest('base64');
}

function hashData(data) {
    const hash = crypto.createHash('sha256');
    hash.update(data);
    return hash.digest('base64');
}

// 解析命令行参数
// -m GET/POST
// [optional] -d 'data'
// -p path (不包括/v1/projects/appid部分)
const args = require('minimist')(process.argv.slice(2), {
    default: {data: '', method: 'GET'},
    string: ['data', 'method', 'path'],
    alias: {d: 'data', m: 'method', p: 'path'},
    boolean: []
});

date = (new Date()).toUTCString();
reqpath = `/v1/projects/${appid}/${args.path}`;
console.log(`path: ${reqpath}`);

reqline = `${args.method} ${reqpath} HTTP/1.1`;
console.log(`request-line: ${reqline}`);

console.log(`data: ${args.data}`);
bodySign = hashData(args.data);
digest = `SHA-256=${bodySign}`;
console.log(`digest: ${digest}`);

signingStr = `host: ${host}\ndate: ${date}\n${reqline}\ndigest: ${digest}`;
console.log(`signingStr: ${signingStr}`);
sign = signData(signingStr);

auth = `hmac username="${username}", `
auth += `algorithm="hmac-sha256", `
auth += `headers="host date request-line digest", `
auth += `signature="${sign}"`;

const options = {
    hostname: host,
    path: `${reqpath}`,
    method: args.method,
    headers: {
        'Date': date,
        'Authorization': auth,
        'Digest': digest,
        'Content-Type': 'application/json'
    }
};

const dataHandler = (res) => {
    if (res.statusCode != 200) {
        console.log(`got unexpected response, statusCode:${
            res.statusCode}, statusMessage:${res.statusMessage}`);
    } else {
        let body = '';
        res.on('data', (chunk) => {
            body += chunk;
        });
        res.on('end', () => {
            console.log(`data: ${body}`);
        });
    }
};

const errHandler = (err) => {
    console.log(`GET ${path} failed, err:${err}`);
};

let req = http.request(options, dataHandler).on('error', errHandler);
req.write(args.data);
req.end();
