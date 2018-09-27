function WsClient(url) {
    this.ws = new WebSocket(url);

    this.eventListener = {};

    this.on = (event, cb) => this.eventListener[event] = cb;

    this.emit = (name, data) => {
        let event = {
            event: name, 
            data: data,
        }
        let rawData = JSON.stringify(event);
        this.ws.send(rawData);
    }

    this.ws.onmessage = (response) => {
        try {
            let data = JSON.parse(response.data);
            if (data) {
                let cb = this.eventListener[data.event]
                if (cb)
                    cb(data.data)
            }
        } catch (e) {
            console.log(e)
        }
    }
}

var ws = new WsClient(
    window.location.href.replace(/((http)|(https)):\/\//gm, 'ws://') + 'ws'
);

// ----------------------------------------------------------------

var $ = (query) => document.querySelector(query);
var guilds;

var tb_token = $('#tb_token');
var btn_submit = $('#btn_submit');
var user_container = $('#user_container');
var lb_name = $('#lb_name');
var lb_tag = $('#lb_tag');
var lb_info = $('#lb_info');
var img_avatar = $('#img_avatar');
var lb_loading = $('#lb_loading');

// ----------------------------------------------------------------

function displayInvalid() {
    tb_token.style.backgroundColor = '#f44336';
    setTimeout(() => {
        tb_token.style.backgroundColor = 'white';
    }, 350);
    user_container.style.display = 'none';
}

// ----------------------------------------------------------------

ws.on('tokenInvalid', () => {
    displayInvalid();
});

ws.on('tokenValid', (data) => {
    tb_token.style.backgroundColor = '#CDDC39';
    setTimeout(() => {
        tb_token.style.backgroundColor = 'white';
    }, 350);

    lb_name.innerText = data.username;
    lb_tag.innerText = data.username + '#' + data.discriminator;
    lb_info.innerHTML = `Running on <b>${data.guilds}</b> servers.`;
    img_avatar.src = data.avatar;

    user_container.style.display = 'flex';
});

ws.on('guildInfo', (data) => {
    lb_loading.style.display = 'none';
    lb_loading.style.animation = '';

    guilds = data;
    window.open('./assets/guildswindow.html', 'Guilds', 'width=800,height=500');
});

user_container.onclick = () => {
    ws.emit('getGuildInfo')
    lb_loading.style.display = 'block';
    lb_loading.style.animation = 'blink 2s ease infinite';
};

btn_submit.onclick = () => {
    let token = tb_token.value;
    if (token.length > 10) {
        ws.emit('checkToken', token);
    } else {
        displayInvalid();
    }
};

