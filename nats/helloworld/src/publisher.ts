import { connect, StringCodec } from "nats";

const main = async () => {
    const nc = await connect({ servers: "demo.nats.io:4222" });

    // create a codec
    const sc = StringCodec();

    nc.publish("hello", sc.encode("Hola"));
    nc.publish("hello", sc.encode("mundo"));

    await nc.drain();
}

main()