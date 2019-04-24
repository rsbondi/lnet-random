const child_process = require("child_process")
const net = require('net')
let spendtimer
let sockpromises = {}
let sockets = {}
let payloadid = 0

const nodes = child_process.execSync("lnet-cli alias")
    .toString()
    .split("\n")
    .reduce((o, n, i) => {
        let node = /alias lcli-([^=]+)="lightning-cli --lightning-dir=([^"]+)"/.exec(n)
        if (node) {
            o.push(node[2])
        }
        return o
    }, [])

spend()

async function spend() {
    const payee = nodes[randInt(nodes.length)]
    let payer
    while (1) {
        payer = nodes[randInt(nodes.length)]
        if (payee != payer) break;
    }
    console.log(`payment from ${payer} to ${payee}`)

    const payeesock = getsocket(payee)
    const payersock = getsocket(payer)

    try {
        const invoice = await postRPC({method: "invoice", params: [randInt(100)*1000, `inv_${payloadid}`, 'send some money']}, payeesock)
        if(invoice && invoice.data && invoice.data.result) await postRPC({method: "pay", params: [invoice.data.result.bolt11]}, payersock)
    } catch(e) {}

    

    if (spendtimer) clearTimeout(spendtimer)
    spendtimer = setTimeout(spend, randInt(10) * 100)
}

function randInt(max) {
    return Math.floor(Math.random() * Math.floor(max));
}

function postRPC(payload, sock) {
    payloadid++
    const promiseFunction = (resolve, reject) => {
        payload.jsonrpc = "2.0"
        payload.params = payload.params || []
        payload.id = payloadid
        sock.write(JSON.stringify(payload))
        sockpromises[payloadid] = {resolve: resolve, reject: reject}
        setTimeout(() => {
            if(sockpromises[payloadid]) {
                sockpromises[payloadid].reject('CONNECTION ERROR')
                sockpromises[payloadid] = undefined
            }
        }, 2000)
    }

    let promise = new Promise(promiseFunction)        
    .catch(e => {
        return console.log(e)
    })

    return promise
}

function getsocket(node) {
    if(sockets[node]) return sockets[node]
    const sock = new net.createConnection(`${node}/lightning-rpc`);
    var _resolve = (key, obj) => {
        sockpromises[key].resolve({data: obj})
        sockpromises[key] = undefined
    }
    sock.on('data', (data) => {
        const response = typeof data == 'string' ? data : data.toString('utf8')
        const obj = JSON.parse(response)
        const key = obj.id
        if(sockpromises && sockpromises[key]) {
            _resolve(key, obj)
        }
    });
    sock.on('error', (derp) => {
        console.log('ERROR:' + derp);
    })
    
    sockets[node] = sock
    return sock
    
}


