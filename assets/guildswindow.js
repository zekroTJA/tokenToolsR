var $ = (query) => document.querySelector(query);
var ws = window.opener.ws;

var t_guilds = $('#t_guilds');
var owner_container = $('#owner_container');
var background_hider = $('#background_hider');
var img_avatar = $('#img_avatar');
var lb_name = $('#lb_name');
var lb_tag = $('#lb_tag');
var lb_id = $('#lb_id');


getUserInfo = (uid) => {
    ws.emit("getUserInfo", uid)
};

hideUserInfo = () => {
    owner_container.style.display = 'none';
    background_hider.style.display = 'none';
    lb_name.innerText = '';
    lb_tag.innerText = '';
    lb_id.innerText = '';
    img_avatar.src = '';
};

ws.on("userInfo", (data) => {
    lb_name.innerText = data.username;
    lb_tag.innerText = data.username + '#' + data.discriminator;
    lb_id.innerText = data.id;
    img_avatar.src = data.avatar;
    owner_container.style.display = 'flex';
    background_hider.style.display = 'block';
});


window.opener.guilds.forEach(g => {
    if (g.name == "")
        return;
    let tr = document.createElement('tr');
    let tdid = document.createElement('td');
    tdid.innerText = g.id;
    let tdname = document.createElement('td');
    tdname.innerText = g.name;
    let tdowner = document.createElement('td');
    tdowner.innerHTML = `<a href="javascript:{}" onclick="getUserInfo('${g.owner}')">${g.owner}</a>`;
    tr.appendChild(tdid);
    tr.appendChild(tdname);
    tr.appendChild(tdowner);
    t_guilds.appendChild(tr);
});

