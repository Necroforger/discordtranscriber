var sounds = {
    micOff: new Audio("/audio/mic-off.wav"),
    micOn: new Audio("/audio/mic-on.wav"),
}
// Make the sounds quieter
sounds.micOn.volume = 0.2;
sounds.micOff.volume = 0.3;

// Create socket
var socket = new ReconnectingWebSocket(`ws://${window.location.host}/websocket/`);

function sendEvent(name, data) {
    socket.send(JSON.stringify({
        Name: name,
        Data: data,
    }));
}


var app = new Vue({
    el: "#app",
    data: {
        avatar: "",
        listening: false,    // If microphone should be listening
        channelValid: false, // If the currently entered channelID is valid

        useVoice: false,     // Use voice synthesis
        useText: true,

        user: undefined,
        channel: undefined,
        guild: undefined,
        voiceChannel: undefined,

        TTS: false,          // Use TTS (text to speech) messages in discord
        channelID: "",
        guildID: "",
        voiceChannelID: "",
        results: [],         // Transcript results
    },
    watch: {
        listening(val) {
            if (val) { // Toggle microphone
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
        channelID(val) { // Request channel information from server
            console.log("CHANNELID: " + val)
            sendEvent("channel", val.trim());
        },
        voiceChannelID(val) { // request voice channel information from server
            console.log("VOICE CHANNEL ID: " + val);
            sendEvent("voice_channel", val.trim());
        },
        voiceChannel(val) {
            this.guildID = val.guild_id;
        },
    },
    computed: {
        hasSpeechRec() {
            return getSpeechRecognition();
        },
        rec() { // Create speech recognition object
            let recognition = getSpeechRecognition();
            let sr = new recognition();

            sr.onresult = (res) => {
                console.log(res);
                let text = res.results[res.results.length - 1][0].transcript.trim();
                this.log(text);
                sendEvent("send", JSON.stringify({
                    ChannelID: this.channelID,
                    GuildID: this.guildID,
                    VoiceChannelID: this.voiceChannelID,
                    Content: text,
                    TTS: this.TTS,
                    Text: (this.useText && (!!this.channel)),
                    Voice: (this.useVoice && (!!this.voiceChannel)),
                }));
            };
            sr.onerror = (err) => {
                console.log(err);
            };
            // Restart recognition after a sentence is recognized.
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
        log(text) {
            this.results.push(text);
            if (this.results.length >= 10) {
                this.results.shift();
            }
        }
    },
});

socket.addEventListener('message', function (event) {
    let ev = JSON.parse(event.data);
    switch (ev.Name) {
        case "valid_channel":
            app.channelValid = JSON.parse(ev.Data); break;
        case "avatar":
            app.avatar = ev.Data; break;
        case "channel":
            app.channel = JSON.parse(ev.Data); break;
        case "voice_channel":
            app.voiceChannel = JSON.parse(ev.Data); break;
        case "user":
            app.user = JSON.parse(ev.Data); break;
        default:
            console.error(`ERROR: Invalid data sent from server: [${ev.Name}]`);
    }
});

function getSpeechRecognition() {
    return window.SpeechRecognition ||
        window.webkitSpeechRecognition ||
        window.mozSpeechRecognition;
}