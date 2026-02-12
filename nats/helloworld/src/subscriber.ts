import { connect, StringCodec } from "nats";

const main = async () => {
    const nc = await connect({ servers: "demo.nats.io:4222" });

    // create a codec
    const sc = StringCodec();
    // create a simple subscriber and iterate over messages
    // matching the subscription
    const sub = nc.subscribe("hello");
    (async () => {
    for await (const m of sub) {
        console.log(`[${sub.getProcessed()}]: ${sc.decode(m.data)}`);
    }
    console.log("subscription closed");
    })();
}

main()