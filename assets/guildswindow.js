var $ = (query) => document.querySelector(query);

var t_guilds = $('#t_guilds');

console.log(window.opener.guilds);

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