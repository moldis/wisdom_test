**Test task for Server Engineer**

Design and implement “Word of Wisdom” tcp server.
• TCP server should be protected from DDOS attacks with the Proof of Work (https://en.wikipedia.org/wiki/Proof_of_work), the challenge-response protocol should be used.
• The choice of the POW algorithm should be explained.
• After Proof Of Work verification, server should send one of the quotes from “word of wisdom” book or any other collection of the quotes.
• Docker file should be provided both for the server and for the client that solves the POW challenge


*I am using simple zk_proof POW algorithm. Not so often used for DDOS, but pretty simple and robust.*


`make build` - will build app

`make server` - will run server, with default cmd flags

`make client` - will connect to default server and try to solve challenge