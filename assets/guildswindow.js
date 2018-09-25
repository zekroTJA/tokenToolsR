var $ = (query) => document.querySelector(query);

var t_guilds = $('#t_guilds');

window.opener.guilds.forEach(g => {
    let tr = document.createElement('tr');
    let tdid = document.createElement('td');
    tdid.innerText = g.id;
    let tdname = document.createElement('td');
    tdname.innerText = g.name;
    let tdowner = document.createElement('td');
    tdowner.innerText = g.owner.name + '#' + g.owner.discriminator;
    tr.appendChild(tdid);
    tr.appendChild(tdname);
    tr.appendChild(tdowner);
    t_guilds.appendChild(tr);
});