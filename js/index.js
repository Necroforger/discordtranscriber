var sounds = {
    micOff: new Audio("/audio/mic-off.wav"),
    micOn: new Audio("/audio/mic-on.wav"),
}
// Make the sounds quieter
sounds.micOn.volume = 0.2;
sounds.micOff.volume = 0.3;

// Create socket
var socket = io()

var app = new Vue({
    el: "#app",
    data: {
        guild: undefined,
        channel: undefined,

        guildID: "",
        channelID: "",
        voiceChannelID: "",
        listening: false,
        results: [],
        avatar: "",
    },
    watch: {
        listening(val) {
            if (val) {
                sounds.micOn.play();
                if (!sounds.micOff.paused) {
                    sounds.micOff.pause();
                    sounds.micOff.currentTime = 0;
                }
                this.startListening();
            } else {
                if (!sounds.micOn.paused) {
                    sounds.micOn.pause();
                    sounds.micOn.currentTime = 0;
                }
                sounds.micOff.play();
                this.stopListening();
            }
        },
        guildID(val) {
            console.log("guildID: " + val);
            socket.emit("guild", val.trim(), (resp) => {
                console.log("updating guild to: " + resp);
                this.guild = resp;
            });
        },
        channelID(val) {
            console.log("channelID: " + val);
            socket.emit("channel", val.trim(), (resp)=> {
                console.log("updating channel to: " + resp);
                this.channel = resp;
            })
        },
    },
    computed: {
        getSpeechRecognition() {
            return getSpeechRecognition();
        },
        rec() {
            let recognition = getSpeechRecognition();
            let sr = new recognition();
            sr.continuous = true;
            sr.onresult = (res) => {
                console.log(res);
                let text = res.results[res.results.length - 1][0].transcript.trim();
                this.results.push(text);
                if (this.results.length >= 10) {
                    this.results.shift();
                }
                socket.emit("send-text", {
                    channelID: this.channelID,
                    content: text,
                });
            };
            sr.onerror = (err) => {
                console.log(err);
            };
            // Begin listening again once a sentence has been recognized
            sr.onend = () => {
                if (app.listening) {
                    sr.start();
                }
            }
            return sr;
        }
    },
    methods: {
        startListening() {
            console.log("Starting speech recognition");
            this.rec.start();
        },
        stopListening() {
            console.log("Stopping speech recognition");
            this.rec.stop();
        },
    },
});

// socket.emit("myAvatar", function (resp) {
//     console.log("Setting avatar to " + resp);
//     app.avatar = resp;
// });

function getSpeechRecognition() {
    return window.SpeechRecognition ||
        window.webkitSpeechRecognition ||
        window.mozSpeechRecognition;
}