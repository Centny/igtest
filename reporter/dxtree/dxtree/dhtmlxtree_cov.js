var g_id = 0;

function cov(ts) {
	var ary = new Array();
	for (var i = 0; i < ts.length; i++) {
		var t = ts[i];
		var a = {};
		a.id = "m_" + (g_id++);
		a.text = t.name;
		a.userdata = [{
			"name": "file",
			"content": t.file
		}, {
			"name": "line",
			"content": t.line
		}, {
			"name": "desc",
			"content": t.desc
		}, {
			"name": "type",
			"content": t.type
		}];
		if (t.subs) {
			a.item = cov(t.subs);
		}
		ary.push(a);
	}
	return ary
}