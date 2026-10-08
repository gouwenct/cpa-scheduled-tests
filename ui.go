package main

func panelHTML() string {
	return `<!doctype html>
<html lang="zh-CN">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>CPA Scheduled Tests</title>
<style>
:root{color-scheme:light dark;--bg:#f5f7fb;--card:#fff;--text:#172033;--muted:#6b7280;--line:#e5e7eb;--accent:#2563eb;--accent2:#1d4ed8;--danger:#dc2626;--ok:#059669;--warn:#d97706;--chip:#eef2ff;--shadow:0 8px 24px rgba(15,23,42,.07)}
@media(prefers-color-scheme:dark){:root{--bg:#0f172a;--card:#111827;--text:#e5e7eb;--muted:#9ca3af;--line:#263244;--accent:#60a5fa;--accent2:#93c5fd;--danger:#f87171;--ok:#34d399;--warn:#fbbf24;--chip:#1e293b;--shadow:none}}
*{box-sizing:border-box}body{margin:0;background:var(--bg);color:var(--text);font:14px/1.5 Inter,ui-sans-serif,system-ui,-apple-system,BlinkMacSystemFont,"Segoe UI",sans-serif}.wrap{max-width:1440px;margin:0 auto;padding:28px}.top{display:flex;justify-content:space-between;gap:20px;align-items:flex-start;margin-bottom:20px}.title h1{font-size:24px;margin:0 0 4px}.title p{margin:0;color:var(--muted)}.keybox{display:flex;gap:8px;align-items:center;flex-wrap:wrap}.keybox input{width:260px}.cards{display:grid;grid-template-columns:repeat(4,minmax(150px,1fr));gap:12px;margin:16px 0 22px}.card,.section{background:var(--card);border:1px solid var(--line);border-radius:12px;box-shadow:var(--shadow)}.card{padding:16px}.metric{font-size:24px;font-weight:700}.label{font-size:12px;color:var(--muted);margin-top:2px}.section{margin:16px 0;padding:18px}.sectionhead{display:flex;justify-content:space-between;gap:12px;align-items:center;margin-bottom:14px}.sectionhead h2{font-size:17px;margin:0}.actions{display:flex;gap:8px;align-items:center;flex-wrap:wrap}button,input,select,textarea{font:inherit}input,select,textarea{background:var(--card);color:var(--text);border:1px solid var(--line);border-radius:8px;padding:9px 10px;outline:none}input:focus,select:focus,textarea:focus{border-color:var(--accent)}button{border:1px solid var(--line);background:var(--card);color:var(--text);border-radius:8px;padding:9px 12px;cursor:pointer}button:hover{border-color:var(--accent)}button.primary{background:var(--accent);border-color:var(--accent);color:white}button.primary:hover{background:var(--accent2)}button.danger{color:var(--danger)}button:disabled{opacity:.5;cursor:not-allowed}.grid2{display:grid;grid-template-columns:1fr 1fr;gap:12px}.grid4{display:grid;grid-template-columns:1.2fr 1.2fr 1fr auto;gap:10px;align-items:end}.field label{display:block;color:var(--muted);font-size:12px;margin:0 0 5px}.field input,.field select,.field textarea{width:100%}.hint{font-size:12px;color:var(--muted);margin-top:6px}.tablewrap{overflow:auto;border:1px solid var(--line);border-radius:10px}table{width:100%;border-collapse:collapse;min-width:980px}th,td{padding:10px 11px;border-bottom:1px solid var(--line);text-align:left;vertical-align:middle}th{font-size:12px;color:var(--muted);font-weight:600;background:color-mix(in srgb,var(--card) 94%,var(--muted) 6%)}tr:last-child td{border-bottom:0}.chip{display:inline-flex;align-items:center;border-radius:999px;padding:2px 8px;font-size:12px;background:var(--chip)}.ok{color:var(--ok)}.bad{color:var(--danger)}.warn{color:var(--warn)}.muted{color:var(--muted)}.nowrap{white-space:nowrap}.ellipsis{max-width:360px;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}.toggle{display:inline-flex;align-items:center;gap:6px}.small{padding:6px 8px;font-size:12px}dialog{width:min(680px,92vw);border:1px solid var(--line);border-radius:12px;background:var(--card);color:var(--text);padding:0;box-shadow:0 24px 80px rgba(0,0,0,.25)}dialog::backdrop{background:rgba(0,0,0,.4)}.modalhead,.modalbody,.modalfoot{padding:16px 18px}.modalhead{border-bottom:1px solid var(--line);font-size:17px;font-weight:700}.modalfoot{border-top:1px solid var(--line);display:flex;justify-content:flex-end;gap:8px}.modalbody{display:grid;gap:12px}.toast{position:fixed;right:20px;bottom:20px;max-width:420px;background:#111827;color:#fff;padding:12px 14px;border-radius:10px;box-shadow:0 10px 30px rgba(0,0,0,.25);display:none;z-index:10}.dangerbox{border:1px solid color-mix(in srgb,var(--danger) 40%,var(--line));background:color-mix(in srgb,var(--danger) 6%,var(--card));padding:12px;border-radius:10px}.tabs{display:flex;gap:8px;margin-bottom:12px}.jobs{display:flex;gap:8px;flex-wrap:wrap}.job{border:1px solid var(--line);border-radius:8px;padding:7px 9px;font-size:12px}.mono{font-family:ui-monospace,SFMono-Regular,Menlo,Consolas,monospace}@media(max-width:900px){.wrap{padding:16px}.top{display:block}.keybox{margin-top:12px}.cards{grid-template-columns:repeat(2,1fr)}.grid2,.grid4{grid-template-columns:1fr}.keybox input{width:100%}}
.github-link{display:inline-flex;align-items:center;gap:6px;margin-top:8px;color:var(--accent);text-decoration:none;border:1px solid var(--line);border-radius:8px;padding:6px 10px}.github-link:hover,.github-link:focus-visible{border-color:var(--accent);text-decoration:underline}.github-link svg{width:18px;height:18px;fill:currentColor}
.plan-pagination{display:flex;justify-content:center;align-items:center;gap:10px;margin-top:12px}.plan-pagination button{padding:6px 10px}
.top{display:block}.keybox{margin-top:12px}.cards{display:flex;gap:8px;margin:0}.card{display:flex;align-items:center;gap:6px;padding:6px 10px;box-shadow:none}.metric{font-size:18px}.label{margin:0;white-space:nowrap}.grid2{grid-template-columns:1fr}.grid4{grid-template-columns:minmax(160px,1fr) minmax(160px,1fr) 90px auto}.sectionhead{display:block}.sectionhead .actions{margin-top:10px}@media(max-width:600px){.grid4{grid-template-columns:1fr}.cards{flex-wrap:wrap}}
</style>
</head>
<body>
<div class="wrap">
  <div class="top">
    <div class="title"><h1>Scheduled Tests</h1><p>CPA 原生的账号 × 模型 × Cron 定时测试。支持单账号立即发送、全部账号强制发送和持久化日志。</p><a class="github-link" href="https://github.com/gouwenct/cpa-scheduled-tests" target="_blank" rel="noopener noreferrer" aria-label="在 GitHub 查看项目并点 Star（新标签页）"><svg viewBox="0 0 24 24" aria-hidden="true"><path d="M12 .5a12 12 0 0 0-3.79 23.39c.6.11.82-.26.82-.58v-2.23c-3.34.73-4.04-1.42-4.04-1.42-.55-1.39-1.33-1.76-1.33-1.76-1.09-.75.08-.73.08-.73 1.2.09 1.84 1.23 1.84 1.23 1.07 1.83 2.81 1.3 3.49.99.11-.77.42-1.3.76-1.6-2.67-.3-5.47-1.33-5.47-5.93 0-1.31.47-2.38 1.24-3.22-.12-.3-.54-1.52.12-3.18 0 0 1.01-.32 3.3 1.23a11.5 11.5 0 0 1 6 0c2.29-1.55 3.3-1.23 3.3-1.23.66 1.66.24 2.88.12 3.18.77.84 1.24 1.91 1.24 3.22 0 4.61-2.8 5.63-5.48 5.93.43.37.82 1.1.82 2.22v3.3c0 .32.22.69.83.57A12 12 0 0 0 12 .5Z"/></svg>GitHub · Star</a></div>
    <div class="keybox"><input id="mgmtKey" type="password" autocomplete="off" placeholder="CPA Management Key"><button onclick="saveKey()">保存本会话</button><button onclick="loadModels(true)">同步模型</button><button class="primary" onclick="loadAll()">刷新</button>
    <div class="cards">
    <div class="card"><div id="mAccounts" class="metric">-</div><div class="label">Codex 账号</div></div>
    <div class="card"><div id="mPlans" class="metric">-</div><div class="label">定时计划</div></div>
    <div class="card"><div id="mEnabled" class="metric">-</div><div class="label">启用计划</div></div>
    </div>
    </div>
  </div>

  <div class="section">
    <div class="sectionhead"><h2>定时测试计划</h2><div class="actions"><button onclick="openBulkPlan()">一键添加全部账号</button><button class="primary" onclick="openPlan()">+ 新建计划</button></div></div>
    <div class="tablewrap"><table><thead><tr><th>启用</th><th>分组</th><th>账号</th><th>Model</th><th>Cron</th><th>时区</th><th>下次执行</th><th>上次结果</th><th>操作</th></tr></thead><tbody id="planRows"></tbody></table></div>
    <div id="planPagination" class="plan-pagination"></div>
    <div class="hint">Cron 为标准 5 段：分钟 小时 日 月 星期。例如 <span class="mono">0 6,11,16,21 * * *</span>。</div>
  </div>

  <div class="section dangerbox">
    <div class="sectionhead"><h2>全部账号立即发送</h2><div class="muted">强制遍历 CPA 中所有 Codex auth 记录，包括 disabled / unavailable；失败账号会写入日志。</div></div>
    <div class="grid4">
      <div class="field"><label>Model</label><select id="allModel"><option>gpt-5.6-luna</option><option>gpt-5.6-sol</option><option>gpt-5.6</option><option>gpt-6-astra</option></select></div>
      <div class="field"><label>Prompt</label><input id="allPrompt" value="你好" placeholder="你好"></div>
      <div class="field"><label>并发</label><input id="bulkConcurrency" type="number" min="1" max="32"></div>
      <button class="primary" onclick="runAllAccounts()">一键向全部账号发送</button>
    </div>
    <div class="hint">每个账号都会发送请求，包括 disabled / unavailable。发送后在日志的“5小时窗口”列查看窗口状态。</div>
  </div>

  <div class="section">
    <div class="sectionhead"><h2>运行任务</h2><div class="actions"><button onclick="loadAll()">刷新状态</button></div></div>
    <div id="jobs" class="jobs"><span class="muted">暂无</span></div>
  </div>

  <div class="section">
    <div class="sectionhead"><h2>日志</h2><div class="actions"><span class="muted">仅保留最近两天</span><select id="logLimit" onchange="loadLogs()"><option>100</option><option selected>200</option><option>500</option><option>1000</option></select><button class="danger" onclick="clearLogs()">清空日志</button></div></div>
    <div class="tablewrap"><table><thead><tr><th>时间</th><th>触发</th><th>账号</th><th>Model</th><th>结果</th><th>5小时窗口</th><th>HTTP</th><th>耗时</th><th>详情</th></tr></thead><tbody id="logRows"></tbody></table></div>
    <div class="hint">窗口状态为本次发送后的快照。只有上游明确返回 300 分钟窗口的用量和有效重置时间，才显示“已开启”或“已开启 · 额度耗尽”；请求成功但缺少窗口数据、以及旧日志均显示“未知”。</div>
  </div>
</div>

<dialog id="planDialog">
  <div class="modalhead" id="planDialogTitle">新建计划</div>
  <div class="modalbody">
    <input id="planId" type="hidden">
    <div class="grid2"><div class="field"><label>分组</label><select id="planGroup"><option>默认</option><option>日常测试</option><option>高频测试</option><option>全账号定时</option></select></div><div class="field"><label>账号</label><select id="planAccount"></select></div></div>
    <div class="grid2"><div class="field"><label>Model</label><select id="planModel"><option>gpt-5.6-luna</option><option>gpt-5.6-sol</option><option>gpt-5.6</option><option>gpt-6-astra</option></select></div><div class="field"><label>时区</label><input id="planTimezone" list="timezoneSuggestions" placeholder="例如 Europe/London"><div class="hint">可从建议中选择，也可填写自己的 IANA 时区。</div></div></div>
    <div class="field"><label>Cron 表达式</label><input id="planCron" class="mono" placeholder="0 6,11,16,21 * * *"><div class="hint">Cron 表达式保留自由输入，默认每天 06:00 / 11:00 / 16:00 / 21:00。</div></div>
    <div class="field"><label>Prompt</label><input id="planPrompt" placeholder="你好"></div>
    <label class="toggle"><input id="planEnabled" type="checkbox" checked> 启用此计划</label>
  </div>
  <div class="modalfoot"><button onclick="document.getElementById('planDialog').close()">取消</button><button class="primary" onclick="savePlan()">保存</button></div>
</dialog>

<dialog id="bulkPlanDialog">
  <div class="modalhead">一键为全部账号创建定时测试</div>
  <div class="modalbody">
    <div class="dangerbox">将为当前 CPA 中的每一个 Codex auth 账号创建一条独立定时计划，包括 disabled / unavailable 账号。已存在完全相同计划的账号会自动跳过，不会重复创建。</div>
    <div class="field"><label>分组</label><select id="bulkPlanGroup"><option>全账号定时</option><option>默认</option><option>日常测试</option><option>高频测试</option></select></div>
    <div class="grid2"><div class="field"><label>Model</label><select id="bulkPlanModel"><option>gpt-5.6-luna</option><option>gpt-5.6-sol</option><option>gpt-5.6</option><option>gpt-6-astra</option></select></div><div class="field"><label>时区</label><input id="bulkPlanTimezone" list="timezoneSuggestions" placeholder="例如 Europe/London"><div class="hint">可从建议中选择，也可填写自己的 IANA 时区。</div></div></div>
    <div class="field"><label>Cron 表达式</label><input id="bulkPlanCron" class="mono" placeholder="0 6,11,16,21 * * *"><div class="hint">同一份时间表会应用到全部账号。例如每天 06:00 / 11:00 / 16:00 / 21:00：<span class="mono">0 6,11,16,21 * * *</span></div></div>
    <div class="field"><label>Prompt</label><input id="bulkPlanPrompt" placeholder="你好"></div>
    <label class="toggle"><input id="bulkPlanEnabled" type="checkbox" checked> 创建后立即启用全部计划</label>
    <div class="hint">创建后仍然是“一账号一计划”，因此你后续可以单独编辑、停用、立即发送或删除某个账号的计划。</div>
  </div>
  <div class="modalfoot"><button onclick="document.getElementById('bulkPlanDialog').close()">取消</button><button class="primary" onclick="createPlansForAll()">创建全部账号计划</button></div>
</dialog>

<div id="toast" class="toast"></div>
<datalist id="timezoneSuggestions"><option value="Asia/Shanghai"><option value="Asia/Taipei"><option value="Asia/Tokyo"><option value="Asia/Singapore"><option value="UTC"><option value="America/Los_Angeles"><option value="America/New_York"><option value="Europe/London"><option value="Europe/Berlin"><option value="Australia/Sydney"></datalist>
<script>
const BASE='/v0/management/plugins/cpa-scheduled-tests';
let STATE={accounts:[],plans:[],settings:{},jobs:[]};
let PLAN_PAGE=1;
const PLAN_PAGE_SIZE=10;
let MODELS=['gpt-5.6-luna','gpt-5.6-sol','gpt-5.6','gpt-6-astra'];
function key(){return document.getElementById('mgmtKey').value.trim()}
function saveKey(){sessionStorage.setItem('cpa-scheduled-tests-key',key());toast('Management Key 已保存到当前标签页会话')}
function headers(json=true){const h={'Authorization':'Bearer '+key()};if(json)h['Content-Type']='application/json';return h}
async function api(path,opt={}){if(!key())throw new Error('请先输入 CPA Management Key');const r=await fetch(BASE+path,{...opt,headers:{...headers(opt.body!==undefined),...(opt.headers||{})}});const t=await r.text();let v;try{v=JSON.parse(t)}catch{v={error:t}}if(!r.ok)throw new Error(v.error||('HTTP '+r.status));return v}
async function hostApi(path,opt={}){if(!key())throw new Error('请先输入 CPA Management Key');const r=await fetch('/v0/management'+path,{...opt,headers:{...headers(opt.body!==undefined),...(opt.headers||{})}});const t=await r.text();let v;try{v=JSON.parse(t)}catch{v={error:t}}if(!r.ok)throw new Error(v.error||('HTTP '+r.status));return v}
async function loadModels(showToast=false){try{const r=await hostApi('/model-definitions/codex',{method:'GET'});const ids=(r.models||[]).map(x=>typeof x==='string'?x:x.id).filter(Boolean);if(ids.length)MODELS=[...new Set(ids)];for(const id of ['planModel','bulkPlanModel','allModel']){const el=document.getElementById(id);const current=el.value;el.innerHTML=MODELS.map(x=>'<option value="'+esc(x)+'">'+esc(x)+'</option>').join('');if(current)setSelectValue(id,current,MODELS[0])}if(showToast)toast('已同步 '+MODELS.length+' 个 Codex Model')}catch(e){if(showToast)toast('同步 Model 失败：'+e.message)}}
function toast(msg){const x=document.getElementById('toast');x.textContent=msg;x.style.display='block';clearTimeout(window.__tt);window.__tt=setTimeout(()=>x.style.display='none',3200)}
function esc(s){return String(s??'').replace(/[&<>'"]/g,c=>({'&':'&amp;','<':'&lt;','>':'&gt;',"'":'&#39;','"':'&quot;'}[c]))}
function fmtTime(s){if(!s||String(s).startsWith('0001-'))return '-';try{const d=new Date(s);const p=n=>String(n).padStart(2,'0');return d.getFullYear()+'.'+p(d.getMonth()+1)+'.'+p(d.getDate())+' '+p(d.getHours())+':'+p(d.getMinutes())+':'+p(d.getSeconds())}catch{return s}}
function statusClass(s){return s==='success'?'ok':s==='limited'?'warn':s?'bad':'muted'}
function fiveHourWindow(x){const s=x.five_hour_status;const label=s==='open'?'已开启':s==='exhausted'?'已开启 · 额度耗尽':'未知';const cls=s==='open'?'ok':s==='exhausted'?'warn':'muted';const detail=[];if(x.five_hour_used_percent!=null&&Number.isFinite(Number(x.five_hour_used_percent)))detail.push('已用 '+Number(x.five_hour_used_percent).toFixed(1)+'%');if(x.five_hour_reset_at)detail.push('重置 '+fmtTime(x.five_hour_reset_at));return '<span class="'+cls+'">'+label+'</span>'+(detail.length?'<br><span class="muted small">'+esc(detail.join(' · '))+'</span>':'')}
function accountLabel(a){const n=a.label||a.email||a.name||a.id||a.auth_index;const flags=[a.disabled?'disabled':'',a.unavailable?'unavailable':''].filter(Boolean);return n+(flags.length?' ['+flags.join(', ')+']':'')}
function fillAccountSelect(id,value){const el=document.getElementById(id);el.innerHTML=STATE.accounts.map(a=>'<option value="'+esc(a.auth_index)+'">'+esc(accountLabel(a))+'</option>').join('');if(value)el.value=value}
function setSelectValue(id,value,fallback){const el=document.getElementById(id);if(value&&!Array.from(el.options).some(o=>o.value===value))el.add(new Option(value,value));el.value=value||fallback}
function browserTimezone(){try{return Intl.DateTimeFormat().resolvedOptions().timeZone||''}catch{return ''}}
function setInputValue(id,value,fallback){document.getElementById(id).value=value||fallback}
function planGroup(p){const raw=String(p?.group||p?.name||'').trim();return raw.split(' · ')[0]||'默认'}
async function loadAll(){try{const s=await api('/state',{method:'GET'});STATE=s;renderState();await loadModels(false);await loadLogs()}catch(e){toast(e.message)}}
function renderState(){document.getElementById('mAccounts').textContent=STATE.accounts.length;document.getElementById('mPlans').textContent=STATE.plans.length;document.getElementById('mEnabled').textContent=STATE.plans.filter(p=>p.enabled).length;fillAccountSelect('planAccount');const set=STATE.settings||{};setSelectValue('allModel',set.default_model||'gpt-5.6-luna','gpt-5.6-luna');document.getElementById('allPrompt').value=document.getElementById('allPrompt').value||'你好';document.getElementById('bulkConcurrency').value=set.bulk_concurrency||4;renderPlans();renderJobs()}
function renderPlans(){const body=document.getElementById('planRows');const totalPages=Math.max(1,Math.ceil(STATE.plans.length/PLAN_PAGE_SIZE));PLAN_PAGE=Math.min(PLAN_PAGE,totalPages);if(!STATE.plans.length){body.innerHTML='<tr><td colspan="9" class="muted">暂无计划。点击“新建计划”。</td></tr>';renderPlanPagination(0,1);return}const start=(PLAN_PAGE-1)*PLAN_PAGE_SIZE;const visible=STATE.plans.slice(start,start+PLAN_PAGE_SIZE);body.innerHTML=visible.map(p=>'<tr><td>'+(p.enabled?'<span class="chip ok">启用</span>':'<span class="chip muted">停用</span>')+'</td><td>'+esc(planGroup(p))+'</td><td class="ellipsis" title="'+esc(p.account_name||p.auth_index)+'">'+esc(p.account_name||p.auth_index)+'</td><td class="mono">'+esc(p.model)+'</td><td class="mono">'+esc(p.cron)+'</td><td>'+esc(p.timezone)+'</td><td class="nowrap">'+fmtTime(p.next_run)+'</td><td><span class="'+statusClass(p.last_status)+'">'+esc(p.last_status||'-')+'</span>'+(p.last_http_status?' · '+p.last_http_status:'')+'</td><td class="nowrap"><button class="small primary" onclick="runPlan(\''+esc(p.id)+'\')">立即发送</button> <button class="small" onclick="togglePlan(\''+esc(p.id)+'\')">'+(p.enabled?'暂停':'恢复')+'</button> <button class="small" onclick="editPlan(\''+esc(p.id)+'\')">编辑</button> <button class="small danger" onclick="deletePlan(\''+esc(p.id)+'\')">删除</button></td></tr>').join('');renderPlanPagination(STATE.plans.length,totalPages)}
function renderPlanPagination(total,totalPages){const el=document.getElementById('planPagination');if(!total){el.innerHTML='';return}el.innerHTML='<button class="small" onclick="changePlanPage(-1)" '+(PLAN_PAGE===1?'disabled':'')+'>上一页</button><span class="muted">第 '+PLAN_PAGE+' / '+totalPages+' 页 · '+total+' 个账号</span><button class="small" onclick="changePlanPage(1)" '+(PLAN_PAGE===totalPages?'disabled':'')+'>下一页</button>'}
function changePlanPage(delta){PLAN_PAGE=Math.max(1,PLAN_PAGE+delta);renderPlans()}
function renderJobs(){const el=document.getElementById('jobs');if(!STATE.jobs.length){el.innerHTML='<span class="muted">暂无</span>';return}el.innerHTML=STATE.jobs.slice(0,10).map(j=>'<div class="job"><b>'+esc(j.kind)+'</b> · '+esc(j.model)+' · '+(j.running?'<span class="warn">运行中</span>':'<span class="muted">完成</span>')+' · '+j.completed+'/'+j.total+' · <span class="ok">'+j.succeeded+' 成功</span> / <span class="warn">'+j.limited+' 限流</span> / <span class="bad">'+j.failed+' 失败</span></div>').join('')}
function openPlan(p=null){document.getElementById('planDialogTitle').textContent=p?'编辑计划':'新建计划';document.getElementById('planId').value=p?.id||'';setSelectValue('planGroup',planGroup(p),'默认');fillAccountSelect('planAccount',p?.auth_index||STATE.accounts[0]?.auth_index||'');setSelectValue('planModel',p?.model||STATE.settings.default_model||'gpt-5.6-luna','gpt-5.6-luna');setInputValue('planTimezone',p?.timezone||browserTimezone()||STATE.settings.default_timezone||'Asia/Shanghai','Asia/Shanghai');document.getElementById('planCron').value=p?.cron||'0 6,11,16,21 * * *';document.getElementById('planPrompt').value=p?.prompt||STATE.settings.default_prompt||'你好';document.getElementById('planEnabled').checked=p?p.enabled:true;document.getElementById('planDialog').showModal()}
function editPlan(id){openPlan(STATE.plans.find(p=>p.id===id))}
function openBulkPlan(){const set=STATE.settings||{};setSelectValue('bulkPlanGroup','全账号定时','全账号定时');setSelectValue('bulkPlanModel',set.default_model||'gpt-5.6-luna','gpt-5.6-luna');setInputValue('bulkPlanTimezone',browserTimezone()||set.default_timezone||'Asia/Shanghai','Asia/Shanghai');document.getElementById('bulkPlanCron').value='0 6,11,16,21 * * *';document.getElementById('bulkPlanPrompt').value=set.default_prompt||'你好';document.getElementById('bulkPlanEnabled').checked=true;document.getElementById('bulkPlanDialog').showModal()}
async function createPlansForAll(){const body={group:document.getElementById('bulkPlanGroup').value,model:document.getElementById('bulkPlanModel').value,timezone:document.getElementById('bulkPlanTimezone').value,cron:document.getElementById('bulkPlanCron').value,prompt:document.getElementById('bulkPlanPrompt').value,enabled:document.getElementById('bulkPlanEnabled').checked};if(!body.model.trim()||!body.cron.trim()){toast('请填写 Model 和 Cron');return}if(!confirm('将为当前 '+STATE.accounts.length+' 个 Codex 账号批量创建定时计划。继续？'))return;try{const r=await api('/plans/bulk-create',{method:'POST',body:JSON.stringify(body)});document.getElementById('bulkPlanDialog').close();toast('批量创建完成：新增 '+r.created+'，跳过重复 '+r.skipped+'，账号总数 '+r.total_accounts);await loadAll()}catch(e){toast(e.message)}}
async function savePlan(){const group=document.getElementById('planGroup').value;const p={id:document.getElementById('planId').value,group,name:group,auth_index:document.getElementById('planAccount').value,model:document.getElementById('planModel').value,timezone:document.getElementById('planTimezone').value,cron:document.getElementById('planCron').value,prompt:document.getElementById('planPrompt').value,enabled:document.getElementById('planEnabled').checked};try{await api('/plans/save',{method:'POST',body:JSON.stringify(p)});document.getElementById('planDialog').close();toast('计划已保存');await loadAll()}catch(e){toast(e.message)}}
async function togglePlan(id){const p=STATE.plans.find(x=>x.id===id);if(!p)return;const body={...p,group:planGroup(p),name:p.name||planGroup(p),enabled:!p.enabled};delete body.next_run;try{await api('/plans/save',{method:'POST',body:JSON.stringify(body)});toast(body.enabled?'计划已恢复':'计划已暂停');await loadAll()}catch(e){toast(e.message)}}
async function deletePlan(id){if(!confirm('确定删除这个计划？'))return;try{await api('/plans/delete',{method:'POST',body:JSON.stringify({id})});toast('计划已删除');await loadAll()}catch(e){toast(e.message)}}
async function runPlan(id){try{const r=await api('/run',{method:'POST',body:JSON.stringify({plan_id:id})});toast('已启动立即发送：'+r.job_id);setTimeout(loadAll,800)}catch(e){toast(e.message)}}
async function runAllAccounts(){const model=document.getElementById('allModel').value.trim();if(!model){toast('请填写 Model');return}if(!confirm('将向 CPA 中所有 Codex 账号强制发送一次真实请求，包括 disabled / unavailable。继续？'))return;try{const set={...STATE.settings,bulk_concurrency:Number(document.getElementById('bulkConcurrency').value||4)};await api('/settings',{method:'POST',body:JSON.stringify(set)});const r=await api('/run-all',{method:'POST',body:JSON.stringify({model,prompt:document.getElementById('allPrompt').value})});toast('全部账号任务已启动：'+r.job_id);setTimeout(loadAll,800)}catch(e){toast(e.message)}}
async function loadLogs(){if(!key())return;try{const n=document.getElementById('logLimit').value||200;const r=await api('/logs?limit='+encodeURIComponent(n),{method:'GET'});const body=document.getElementById('logRows');const logs=r.logs||[];if(!logs.length){body.innerHTML='<tr><td colspan="9" class="muted">暂无日志</td></tr>';return}body.innerHTML=logs.map(x=>'<tr><td class="nowrap">'+fmtTime(x.at)+'</td><td>'+esc(x.trigger)+'</td><td class="ellipsis" title="'+esc(x.account_name)+'">'+esc(x.account_name)+(x.disabled?' <span class="chip muted">disabled</span>':'')+(x.unavailable?' <span class="chip warn">unavailable</span>':'')+'</td><td class="mono">'+esc(x.model)+'</td><td class="'+statusClass(x.status)+'">'+esc(x.status)+'</td><td>'+fiveHourWindow(x)+'</td><td>'+(x.http_status||'-')+'</td><td>'+x.latency_ms+' ms</td><td class="ellipsis" title="'+esc(x.error||'')+'">'+esc(x.error||'-')+'</td></tr>').join('')}catch(e){toast(e.message)}}
async function clearLogs(){if(!confirm('确定清空插件运行日志？'))return;try{await api('/logs/clear',{method:'POST',body:'{}'});toast('日志已清空');await loadLogs()}catch(e){toast(e.message)}}
window.addEventListener('load',()=>{const k=sessionStorage.getItem('cpa-scheduled-tests-key')||'';document.getElementById('mgmtKey').value=k;if(k)loadAll()});
setInterval(()=>{if(key()&&STATE.jobs?.some(j=>j.running))loadAll()},2500);
</script>
</body></html>`
}
