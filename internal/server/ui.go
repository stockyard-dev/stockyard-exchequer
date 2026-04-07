package server

import "net/http"

func (s *Server) dashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(dashHTML))
}

const dashHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width,initial-scale=1.0">
<title>Exchequer</title>
<link href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;500;700&display=swap" rel="stylesheet">
<style>
:root{--bg:#1a1410;--bg2:#241e18;--bg3:#2e261e;--rust:#e8753a;--leather:#a0845c;--cream:#f0e6d3;--cd:#bfb5a3;--cm:#7a7060;--gold:#d4a843;--green:#4a9e5c;--red:#c94444;--orange:#d4843a;--blue:#5b8dd9;--mono:'JetBrains Mono',monospace}
*{margin:0;padding:0;box-sizing:border-box}
body{background:var(--bg);color:var(--cream);font-family:var(--mono);line-height:1.5;font-size:13px}
.hdr{padding:.8rem 1.5rem;border-bottom:1px solid var(--bg3);display:flex;justify-content:space-between;align-items:center;gap:1rem;flex-wrap:wrap}
.hdr h1{font-size:.9rem;letter-spacing:2px}
.hdr h1 span{color:var(--rust)}
.main{padding:1.2rem 1.5rem;max-width:1000px;margin:0 auto}
.stats{display:grid;grid-template-columns:repeat(4,1fr);gap:.5rem;margin-bottom:1rem}
.st{background:var(--bg2);border:1px solid var(--bg3);padding:.7rem;text-align:center}
.st-v{font-size:1.2rem;font-weight:700;color:var(--gold)}
.st-v.green{color:var(--green)}
.st-v.warn{color:var(--orange)}
.st-v.red{color:var(--red)}
.st-l{font-size:.5rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-top:.2rem}
.toolbar{display:flex;gap:.5rem;margin-bottom:1rem;flex-wrap:wrap;align-items:center}
.search{flex:1;min-width:180px;padding:.4rem .6rem;background:var(--bg2);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.7rem}
.search:focus{outline:none;border-color:var(--leather)}
.filter-sel{padding:.4rem .5rem;background:var(--bg2);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.65rem}
.item{background:var(--bg2);border:1px solid var(--bg3);padding:.9rem 1rem;margin-bottom:.5rem;transition:border-color .15s}
.item:hover{border-color:var(--leather)}
.item.over{border-left:3px solid var(--red)}
.item.warn{border-left:3px solid var(--orange)}
.item.healthy{border-left:3px solid var(--green)}
.item-top{display:flex;justify-content:space-between;align-items:flex-start;gap:.5rem;margin-bottom:.5rem}
.item-title-block{flex:1;min-width:0}
.item-title{font-size:.85rem;font-weight:700;color:var(--cream)}
.item-sub{font-size:.6rem;color:var(--cm);margin-top:.15rem;display:flex;gap:.6rem;flex-wrap:wrap;align-items:center}
.item-actions{display:flex;gap:.3rem;flex-shrink:0}
.budget-row{display:flex;justify-content:space-between;align-items:center;font-size:.7rem;margin-bottom:.3rem}
.budget-spent{color:var(--cream);font-weight:700}
.budget-spent.over{color:var(--red)}
.budget-allocated{color:var(--cm)}
.budget-bar{height:8px;background:var(--bg3);overflow:hidden;position:relative}
.budget-fill{height:100%;transition:width .3s,background-color .3s}
.budget-fill.healthy{background:var(--green)}
.budget-fill.warn{background:var(--orange)}
.budget-fill.over{background:var(--red)}
.budget-info{display:flex;justify-content:space-between;font-size:.55rem;color:var(--cm);margin-top:.25rem}
.budget-info .pct{font-weight:700}
.budget-info .pct.over{color:var(--red)}
.budget-info .pct.warn{color:var(--orange)}
.budget-info .remaining{color:var(--cd)}
.budget-info .remaining.negative{color:var(--red);font-weight:700}
.item-meta{font-size:.55rem;color:var(--cm);margin-top:.5rem;display:flex;gap:.6rem;flex-wrap:wrap}
.item-extra{font-size:.55rem;color:var(--cd);margin-top:.4rem;padding-top:.3rem;border-top:1px dashed var(--bg3);display:flex;flex-direction:column;gap:.15rem}
.item-extra-row{display:flex;gap:.4rem}
.item-extra-label{color:var(--cm);text-transform:uppercase;letter-spacing:.5px;min-width:90px}
.item-extra-val{color:var(--cream)}
.badge{font-size:.5rem;padding:.12rem .35rem;text-transform:uppercase;letter-spacing:1px;border:1px solid var(--bg3);color:var(--cm);font-weight:700}
.badge.cat{border-color:var(--leather);color:var(--leather)}
.badge.period{border-color:var(--blue);color:var(--blue)}
.btn{font-family:var(--mono);font-size:.6rem;padding:.25rem .5rem;cursor:pointer;border:1px solid var(--bg3);background:var(--bg);color:var(--cd);transition:.15s}
.btn:hover{border-color:var(--leather);color:var(--cream)}
.btn-p{background:var(--rust);border-color:var(--rust);color:#fff}
.btn-p:hover{opacity:.85;color:#fff}
.btn-sm{font-size:.55rem;padding:.2rem .4rem}
.modal-bg{display:none;position:fixed;inset:0;background:rgba(0,0,0,.65);z-index:100;align-items:center;justify-content:center}
.modal-bg.open{display:flex}
.modal{background:var(--bg2);border:1px solid var(--bg3);padding:1.5rem;width:520px;max-width:92vw;max-height:90vh;overflow-y:auto}
.modal h2{font-size:.8rem;margin-bottom:1rem;color:var(--rust);letter-spacing:1px}
.fr{margin-bottom:.6rem}
.fr label{display:block;font-size:.55rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-bottom:.2rem}
.fr input,.fr select,.fr textarea{width:100%;padding:.4rem .5rem;background:var(--bg);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.7rem}
.fr input:focus,.fr select:focus,.fr textarea:focus{outline:none;border-color:var(--leather)}
.row2{display:grid;grid-template-columns:1fr 1fr;gap:.5rem}
.fr-section{margin-top:1rem;padding-top:.8rem;border-top:1px solid var(--bg3)}
.fr-section-label{font-size:.55rem;color:var(--rust);text-transform:uppercase;letter-spacing:1px;margin-bottom:.5rem}
.acts{display:flex;gap:.4rem;justify-content:flex-end;margin-top:1rem}
.acts .btn-del{margin-right:auto;color:var(--red);border-color:#3a1a1a}
.acts .btn-del:hover{border-color:var(--red);color:var(--red)}
.empty{text-align:center;padding:3rem;color:var(--cm);font-style:italic;font-size:.85rem}
@media(max-width:600px){.stats{grid-template-columns:repeat(2,1fr)}}
</style>
</head>
<body>

<div class="hdr">
<h1 id="dash-title"><span>&#9670;</span> EXCHEQUER</h1>
<button class="btn btn-p" onclick="openForm()">+ Add Budget</button>
</div>

<div class="main">
<div class="stats" id="stats"></div>
<div class="toolbar">
<input class="search" id="search" placeholder="Search budgets..." oninput="debouncedRender()">
<select class="filter-sel" id="period-filter" onchange="render()">
<option value="">All Periods</option>
<option value="monthly">Monthly</option>
<option value="quarterly">Quarterly</option>
<option value="yearly">Yearly</option>
<option value="project">Project</option>
</select>
<select class="filter-sel" id="category-filter" onchange="render()">
<option value="">All Categories</option>
</select>
</div>
<div id="list"></div>
</div>

<div class="modal-bg" id="mbg" onclick="if(event.target===this)closeModal()">
<div class="modal" id="mdl"></div>
</div>

<script>
var A='/api';
var RESOURCE='budgets';

var fields=[
{name:'name',label:'Budget Name',type:'text',required:true},
{name:'category',label:'Category',type:'select_or_text',options:[]},
{name:'allocated',label:'Allocated ($)',type:'money',required:true},
{name:'spent',label:'Spent ($)',type:'money'},
{name:'period',label:'Period',type:'select',options:['monthly','quarterly','yearly','project']},
{name:'start_date',label:'Start Date',type:'date'},
{name:'end_date',label:'End Date',type:'date'},
{name:'notes',label:'Notes',type:'textarea'}
];

var budgets=[],budgetExtras={},editId=null,searchTimer=null;

// ─── Money helpers ────────────────────────────────────────────────

function fmtMoney(cents){
if(!cents&&cents!==0)return'$0';
var dollars=cents/100;
var negative=dollars<0;
var abs=Math.abs(dollars);
var str;
if(abs>=1000000)str='$'+(abs/1000000).toFixed(1)+'M';
else if(abs>=1000)str='$'+(abs/1000).toFixed(1)+'k';
else str='$'+abs.toFixed(0);
return negative?'-'+str:str;
}

function fmtMoneyFull(cents){
if(!cents&&cents!==0)return'$0';
var negative=cents<0;
var abs=Math.abs(cents)/100;
var str='$'+abs.toLocaleString('en-US',{minimumFractionDigits:0,maximumFractionDigits:0});
return negative?'-'+str:str;
}

function parseMoney(str){
if(!str)return 0;
var n=parseFloat(String(str).replace(/[^\d.-]/g,''));
if(isNaN(n))return 0;
return Math.round(n*100);
}

function fieldByName(n){
for(var i=0;i<fields.length;i++)if(fields[i].name===n)return fields[i];
return null;
}

function debouncedRender(){
clearTimeout(searchTimer);
searchTimer=setTimeout(render,200);
}

// ─── Loading ──────────────────────────────────────────────────────

async function load(){
try{
var resps=await Promise.all([
fetch(A+'/budgets').then(function(r){return r.json()}),
fetch(A+'/stats').then(function(r){return r.json()})
]);
budgets=resps[0].budgets||[];
renderStats(resps[1]||{});

try{
var ex=await fetch(A+'/extras/'+RESOURCE).then(function(r){return r.json()});
budgetExtras=ex||{};
budgets.forEach(function(b){
var x=budgetExtras[b.id];
if(!x)return;
Object.keys(x).forEach(function(k){if(b[k]===undefined)b[k]=x[k]});
});
}catch(e){budgetExtras={}}

populateCategoryFilter();
}catch(e){
console.error('load failed',e);
budgets=[];
}
render();
}

function populateCategoryFilter(){
var sel=document.getElementById('category-filter');
if(!sel)return;
var current=sel.value;
var seen={};
var cats=[];
budgets.forEach(function(b){
if(b.category&&!seen[b.category]){seen[b.category]=true;cats.push(b.category)}
});
cats.sort();
sel.innerHTML='<option value="">All Categories</option>'+cats.map(function(c){return'<option value="'+esc(c)+'"'+(c===current?' selected':'')+'>'+esc(c)+'</option>'}).join('');
}

function renderStats(s){
var total=s.total||0;
var allocated=s.total_allocated||0;
var spent=s.total_spent||0;
var remaining=s.remaining||0;
var overBudget=s.over_budget||0;
var remainingClass=remaining<0?'red':'green';
document.getElementById('stats').innerHTML=
'<div class="st"><div class="st-v">'+fmtMoney(allocated)+'</div><div class="st-l">Allocated</div></div>'+
'<div class="st"><div class="st-v warn">'+fmtMoney(spent)+'</div><div class="st-l">Spent</div></div>'+
'<div class="st"><div class="st-v '+remainingClass+'">'+fmtMoney(remaining)+'</div><div class="st-l">Remaining</div></div>'+
'<div class="st"><div class="st-v '+(overBudget>0?'red':'')+'">'+overBudget+'</div><div class="st-l">Over Budget</div></div>';
}

function render(){
var q=(document.getElementById('search').value||'').toLowerCase();
var pf=document.getElementById('period-filter').value;
var cf=document.getElementById('category-filter').value;

var f=budgets;
if(q)f=f.filter(function(b){
return(b.name||'').toLowerCase().includes(q)||
(b.category||'').toLowerCase().includes(q)||
(b.notes||'').toLowerCase().includes(q);
});
if(pf)f=f.filter(function(b){return b.period===pf});
if(cf)f=f.filter(function(b){return b.category===cf});

if(!f.length){
var msg=window._emptyMsg||'No budgets. Add your first one to start tracking.';
document.getElementById('list').innerHTML='<div class="empty">'+esc(msg)+'</div>';
return;
}

var h='';
f.forEach(function(b){h+=itemHTML(b)});
document.getElementById('list').innerHTML=h;
}

function itemHTML(b){
var allocated=b.allocated||0;
var spent=b.spent||0;
var pct=allocated>0?Math.round((spent/allocated)*100):0;
var displayPct=Math.min(pct,100);
var remaining=allocated-spent;
var status='healthy';
if(pct>=100)status='over';
else if(pct>=80)status='warn';

var h='<div class="item '+status+'">';
h+='<div class="item-top"><div class="item-title-block">';
h+='<div class="item-title">'+esc(b.name)+'</div>';
h+='<div class="item-sub">';
if(b.category)h+='<span class="badge cat">'+esc(b.category)+'</span>';
if(b.period)h+='<span class="badge period">'+esc(b.period)+'</span>';
h+='</div>';
h+='</div>';
h+='<div class="item-actions"><button class="btn btn-sm" onclick="openEdit(\''+esc(b.id)+'\')">Edit</button></div>';
h+='</div>';

// Budget bar
h+='<div class="budget-row">';
h+='<span class="budget-spent '+(pct>=100?'over':'')+'">'+fmtMoneyFull(spent)+'</span>';
h+='<span class="budget-allocated">of '+fmtMoneyFull(allocated)+'</span>';
h+='</div>';
h+='<div class="budget-bar"><div class="budget-fill '+status+'" style="width:'+displayPct+'%"></div></div>';
h+='<div class="budget-info">';
h+='<span class="pct '+status+'">'+pct+'%</span>';
h+='<span class="remaining'+(remaining<0?' negative':'')+'">'+(remaining<0?'over by ':'')+fmtMoneyFull(Math.abs(remaining))+(remaining<0?'':' left')+'</span>';
h+='</div>';

if(b.start_date||b.end_date||b.notes){
h+='<div class="item-meta">';
if(b.start_date)h+='<span>Start: '+esc(b.start_date)+'</span>';
if(b.end_date)h+='<span>End: '+esc(b.end_date)+'</span>';
h+='</div>';
}

// Custom fields
var customRows='';
fields.forEach(function(f){
if(!f.isCustom)return;
var v=b[f.name];
if(v===undefined||v===null||v==='')return;
customRows+='<div class="item-extra-row">';
customRows+='<span class="item-extra-label">'+esc(f.label)+'</span>';
customRows+='<span class="item-extra-val">'+esc(String(v))+'</span>';
customRows+='</div>';
});
if(customRows)h+='<div class="item-extra">'+customRows+'</div>';

h+='</div>';
return h;
}

// ─── Modal ────────────────────────────────────────────────────────

function fieldHTML(f,value){
var v=value;
if(v===undefined||v===null)v='';
var req=f.required?' *':'';
var ph=f.placeholder?(' placeholder="'+esc(f.placeholder)+'"'):'';
var h='<div class="fr"><label>'+esc(f.label)+req+'</label>';

if(f.type==='select'){
h+='<select id="f-'+f.name+'">';
if(!f.required)h+='<option value="">Select...</option>';
(f.options||[]).forEach(function(o){
var sel=(String(v)===String(o))?' selected':'';
var disp=String(o).charAt(0).toUpperCase()+String(o).slice(1).replace(/_/g,' ');
h+='<option value="'+esc(String(o))+'"'+sel+'>'+esc(disp)+'</option>';
});
h+='</select>';
}else if(f.type==='select_or_text'){
h+='<input list="dl-'+f.name+'" type="text" id="f-'+f.name+'" value="'+esc(String(v))+'"'+ph+'>';
h+='<datalist id="dl-'+f.name+'">';
var opts=(f.options||[]).slice();
budgets.forEach(function(bd){
if(bd.category&&opts.indexOf(bd.category)===-1)opts.push(bd.category);
});
opts.forEach(function(o){h+='<option value="'+esc(String(o))+'">'});
h+='</datalist>';
}else if(f.type==='textarea'){
h+='<textarea id="f-'+f.name+'" rows="2"'+ph+'>'+esc(String(v))+'</textarea>';
}else if(f.type==='money'){
var dollars=v?(v/100).toFixed(2):'';
h+='<input type="text" id="f-'+f.name+'" value="'+esc(dollars)+'" placeholder="0.00">';
}else if(f.type==='number'||f.type==='integer'){
h+='<input type="number" id="f-'+f.name+'" value="'+esc(String(v))+'"'+ph+'>';
}else{
var inputType=f.type||'text';
h+='<input type="'+esc(inputType)+'" id="f-'+f.name+'" value="'+esc(String(v))+'"'+ph+'>';
}
h+='</div>';
return h;
}

function formHTML(budget){
var b=budget||{};
var isEdit=!!budget;
var h='<h2>'+(isEdit?'EDIT BUDGET':'NEW BUDGET')+'</h2>';

h+=fieldHTML(fieldByName('name'),b.name);
h+='<div class="row2">'+fieldHTML(fieldByName('category'),b.category)+fieldHTML(fieldByName('period'),b.period||'monthly')+'</div>';
h+='<div class="row2">'+fieldHTML(fieldByName('allocated'),b.allocated)+fieldHTML(fieldByName('spent'),b.spent)+'</div>';
h+='<div class="row2">'+fieldHTML(fieldByName('start_date'),b.start_date)+fieldHTML(fieldByName('end_date'),b.end_date)+'</div>';
h+=fieldHTML(fieldByName('notes'),b.notes);

var customFields=fields.filter(function(f){return f.isCustom});
if(customFields.length){
var label=window._customSectionLabel||'Additional Details';
h+='<div class="fr-section"><div class="fr-section-label">'+esc(label)+'</div>';
customFields.forEach(function(f){h+=fieldHTML(f,b[f.name])});
h+='</div>';
}

h+='<div class="acts">';
if(isEdit){
h+='<button class="btn btn-del" onclick="delBudget()">Delete</button>';
}
h+='<button class="btn" onclick="closeModal()">Cancel</button>';
h+='<button class="btn btn-p" onclick="submit()">'+(isEdit?'Save':'Add')+'</button>';
h+='</div>';
return h;
}

function openForm(){
editId=null;
document.getElementById('mdl').innerHTML=formHTML();
document.getElementById('mbg').classList.add('open');
var n=document.getElementById('f-name');
if(n)n.focus();
}

function openEdit(id){
var b=null;
for(var i=0;i<budgets.length;i++){if(budgets[i].id===id){b=budgets[i];break}}
if(!b)return;
editId=id;
document.getElementById('mdl').innerHTML=formHTML(b);
document.getElementById('mbg').classList.add('open');
}

function closeModal(){
document.getElementById('mbg').classList.remove('open');
editId=null;
}

async function submit(){
var nameEl=document.getElementById('f-name');
if(!nameEl||!nameEl.value.trim()){alert('Budget name is required');return}

var body={};
var extras={};
fields.forEach(function(f){
var el=document.getElementById('f-'+f.name);
if(!el)return;
var val;
if(f.type==='money')val=parseMoney(el.value);
else if(f.type==='number'||f.type==='integer')val=parseFloat(el.value)||0;
else val=el.value.trim();
if(f.isCustom)extras[f.name]=val;
else body[f.name]=val;
});

var savedId=editId;
try{
if(editId){
var r1=await fetch(A+'/budgets/'+editId,{method:'PUT',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});
if(!r1.ok){var e1=await r1.json().catch(function(){return{}});alert(e1.error||'Save failed');return}
}else{
var r2=await fetch(A+'/budgets',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});
if(!r2.ok){var e2=await r2.json().catch(function(){return{}});alert(e2.error||'Add failed');return}
var created=await r2.json();
savedId=created.id;
}
if(savedId&&Object.keys(extras).length){
await fetch(A+'/extras/'+RESOURCE+'/'+savedId,{method:'PUT',headers:{'Content-Type':'application/json'},body:JSON.stringify(extras)}).catch(function(){});
}
}catch(e){
alert('Network error: '+e.message);
return;
}
closeModal();
load();
}

async function delBudget(){
if(!editId)return;
if(!confirm('Delete this budget?'))return;
await fetch(A+'/budgets/'+editId,{method:'DELETE'});
closeModal();
load();
}

function esc(s){
if(s===undefined||s===null)return'';
var d=document.createElement('div');
d.textContent=String(s);
return d.innerHTML;
}

document.addEventListener('keydown',function(e){if(e.key==='Escape')closeModal()});

// ─── Personalization ──────────────────────────────────────────────

(function loadPersonalization(){
fetch('/api/config').then(function(r){return r.json()}).then(function(cfg){
if(!cfg||typeof cfg!=='object')return;

if(cfg.dashboard_title){
var h1=document.getElementById('dash-title');
if(h1)h1.innerHTML='<span>&#9670;</span> '+esc(cfg.dashboard_title);
document.title=cfg.dashboard_title;
}

if(cfg.empty_state_message)window._emptyMsg=cfg.empty_state_message;
if(cfg.primary_label)window._customSectionLabel=cfg.primary_label+' Details';

if(Array.isArray(cfg.categories)){
var catField=fieldByName('category');
if(catField)catField.options=cfg.categories;
}

if(Array.isArray(cfg.custom_fields)){
cfg.custom_fields.forEach(function(cf){
if(!cf||!cf.name||!cf.label)return;
if(fieldByName(cf.name))return;
fields.push({
name:cf.name,
label:cf.label,
type:cf.type||'text',
options:cf.options||[],
isCustom:true
});
});
}
}).catch(function(){
}).finally(function(){
load();
});
})();
</script>
</body>
</html>`
