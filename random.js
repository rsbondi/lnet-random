(async function() {
    const fs = require('fs')

    const args = process.argv.reduce((o, a) => {
        if (~a.indexOf('=')) {
            const kv = a.split('=')
            o[kv[0]] = kv[1]
        }
        return o
    }, {})
    
    let nodes = []

    try {
        let users = await get(args.count)
        users.results.forEach(u => {
            nodes.push(u.name.last)
        })
    } catch (e) {
        for (let i = 0; i < args.count; i++) {
            nodes.push(Math.random().toString(36).substring(2, 15) + Math.random().toString(36).substring(2, 15))
        }
    }
    
    let graph = ''
    let cons = {}
    for (let i = 0; i < nodes.length; i++) {
        const node = nodes[i]
        const peers = Math.floor(Math.random() * (args.maxchannels || 3)) + 1
        cons[node] = []
        for (let p = 0; p < peers; p++) {
            const d = Math.floor(Math.random() * args.count)
            if (d === i) { p--; continue } // con't connect to self, try again
            const src = nodes[d]
            if (~cons[node].indexOf(src)) { p--; continue } // only one connection
            if (cons[src] && ~cons[src].indexOf(node)) { p--; continue } // the other way
            cons[node].push(src)
            const cap = Math.floor(Math.random() * (10000000 - 500000 + 1) + 500000)
            graph += `  "${src}" -- "${node}" [capacity="${cap}"];\n`
        }
    }
    
    fs.writeFileSync(args.out, Buffer.from(`graph g {\n${graph}}`))
})()

function get(n) {
    const https = require('https')
    return new Promise((resolve, reject) => {

        const req = https.request({
            host:'randomuser.me', 
            method: 'GET', 
            port: '443',
            path: `/api/?results=${n}&inc=name`}, (res) => {
            res.setEncoding('utf8');
            let response = ''
            res.on('data', (chunk) => {
              response += chunk
            });
            res.on('end', () => {
              resolve(JSON.parse(response))
            });
          });
      
          req.on('error', (e) => {
            reject(`problem with request: ${e.message}`);
          });
      
          req.end();  

    })
}
