(this["webpackJsonpreact-table-demo"]=this["webpackJsonpreact-table-demo"]||[]).push([[0],[,,,function(e,t,a){e.exports={root:"Table_root__3x77O",fill:"Table_fill__3Frg4",Table:"Table_Table__gTzsf",Row:"Table_Row__29GKa",BodyCell:"Table_BodyCell__2gQ09",HeaderCell:"Table_HeaderCell__37MIo",ResizeHandle:"Table_ResizeHandle__1LDIp",ResizeHandleActive:"Table_ResizeHandleActive__-UEDL"}},,,,function(e,t,a){e.exports={ColumnSelector:"ColumnSelector_ColumnSelector__1LtuL",Label:"ColumnSelector_Label__Znnco",Checkbox:"ColumnSelector_Checkbox__3CcNf"}},,,,function(e,t,a){"use strict";(function(e){a.d(t,"a",(function(){return d}));var n=a(14),c=a(4),l=a(1),r=(a(21),a(0));function s(e){return e[Math.floor(Math.random()*e.length)]}function o(e,t){var a=s(["200 OK","301 Moved Permanently","404 Not Found","418 I'm a teapot","501 Not Implemented"]);return{timestamp:e+Math.floor(Math.random()*(t-e)),latencyMs:5+Math.floor(150*Math.random()),endpoint:"/user/"+s(["bendy-badger","happy-hippo","giant-ape","grumpy-groundhog","phlegmatic-pheasant"]),status:a}}var i=[{Header:"Timestamp",Cell:function(e){var t=e.value;return Object(r.jsx)("span",{className:"Cell-Timestamp",children:Object(r.jsx)("span",{children:new Date(t).toLocaleString()})})},accessor:"timestamp"},{Header:"Latency",Cell:function(e){var t=e.value,a="bad";return t<=50?a="good":t<=100&&(a="weak"),Object(r.jsx)("span",{className:"Cell-Latency ".concat(a),children:t})},accessor:"latencyMs"},{Header:"Endpoint",Cell:function(e){var t=e.value;return Object(r.jsx)("span",{className:"Cell-Endpoint",children:t})},accessor:"endpoint"},{Header:"Status",Cell:function(e){var t=e.value,a=+t.split(" ")[0],n=500;return a<300?n=200:a<400?n=300:a<500&&(n=400),Object(r.jsx)("span",{className:"Cell-StatusCode range-".concat(n),children:t})},accessor:"status"}];function d(){var t=arguments.length>0&&void 0!==arguments[0]?arguments[0]:5,a=arguments.length>1&&void 0!==arguments[1]?arguments[1]:1e3,r=arguments.length>2&&void 0!==arguments[2]?arguments[2]:100,s=l.useState(Date.now()-6048e5),d=Object(c.a)(s,2),u=d[0],b=d[1],j=l.useState([]),h=Object(c.a)(j,2),m=h[0],p=h[1],O=l.useCallback((function(){if(!(m.length>=r)){var e=Array(t).fill(0).map((function(e,t){return o(u+1e3*t,u+1999*t)}));p([].concat(Object(n.a)(m),[e])),b(u+2e3*t)}}),[m,r,u,t]);return l.useEffect((function(){var t=e.setInterval(O,a);return function(){e.clearInterval(t)}}),[a,O]),l.useMemo((function(){return{columns:i,data:m.flat()}}),[m])}}).call(this,a(20))},function(e,t,a){e.exports={Filter:"Filter_Filter__3qIm_"}},,,,,,,function(e,t,a){},,function(e,t,a){},,,,,function(e,t,a){},function(e,t,a){"use strict";a.r(t);var n=a(1),c=a.n(n),l=a(10),r=a.n(l),s=(a(19),a(11)),o=a(2),i=a(4),d=a(3),u=a.n(d),b=a(12),j=a.n(b),h=a(0);function m(e){var t=e.onChange,a=n.useState(""),c=Object(i.a)(a,2),l=c[0],r=c[1],s=n.useCallback((function(e){var a=e.target.value.trim();r(a),t(a)}),[t]);return Object(h.jsx)("div",{className:j.a.Filter,children:Object(h.jsx)("input",{type:"text",value:l,placeholder:"Search rows...",onChange:s})})}var p=a(7),O=a.n(p);function f(e){var t=e.columns;return Object(h.jsxs)("div",{className:O.a.ColumnSelector,children:[Object(h.jsx)("div",{className:O.a.Label,children:"Show Columns:"}),Object(h.jsx)("div",{className:O.a.Checkboxes,children:t.map((function(e){return Object(h.jsx)("div",{className:O.a.Checkbox,children:Object(h.jsxs)("label",{children:[Object(h.jsx)("input",Object(o.a)({type:"checkbox"},e.getToggleHiddenProps()))," ".concat(e.Header)]})},e.id)}))})]})}var v=a(5),g=a(13);function x(e){var t,a=e.data,c=a.columns,l=a.data,r=Object(v.useTable)({columns:c,data:l},v.useFlexLayout,v.useGlobalFilter,v.useSortBy,v.useResizeColumns),s=r.getTableProps,d=r.getTableBodyProps,b=r.headerGroups,j=r.rows,p=r.allColumns,O=r.prepareRow,x=r.setGlobalFilter,C=n.useMemo((function(){var e=document.createElement("div");e.setAttribute("style","width: 100vw; height: 100vh; overflow: scroll; position: absolute; top: -100vh;"),document.body.appendChild(e);var t=e.offsetWidth-e.clientWidth,a=e.offsetHeight-e.clientHeight;return document.body.removeChild(e),{width:t,height:a}}),[]).width,_=n.useState(null),y=Object(i.a)(_,2),T=y[0],w=y[1],N=n.useCallback((function(e){return w(e)}),[]),H=(t=T,n.useMemo((function(){return t?{width:t.clientWidth,height:t.clientHeight}:{width:0,height:0}}),[t])).height;return Object(h.jsxs)("div",{className:u.a.root,children:[Object(h.jsxs)("header",{children:[Object(h.jsx)(f,{columns:p}),Object(h.jsx)(m,{onChange:x})]}),Object(h.jsx)("div",{className:u.a.fill,ref:N,children:Object(h.jsxs)("div",Object(o.a)(Object(o.a)({},s()),{},{className:u.a.Table,children:[Object(h.jsx)("div",{className:u.a.TableHead,children:b.map((function(e,t){return Object(h.jsx)("div",Object(o.a)(Object(o.a)({className:u.a.Row},e.getHeaderGroupProps()),{},{children:e.headers.map((function(a,n){return Object(h.jsxs)("div",Object(o.a)(Object(o.a)({className:u.a.HeaderCell},a.getHeaderProps({style:t===b.length-1&&n===e.headers.length-1?{marginRight:C}:{}})),{},{children:[Object(h.jsxs)("div",Object(o.a)(Object(o.a)({},a.getSortByToggleProps()),{},{children:[a.render("Header"),Object(h.jsx)("span",{children:a.isSorted?a.isSortedDesc?" \ud83d\udd3d":" \ud83d\udd3c":""})]})),Object(h.jsx)("div",Object(o.a)(Object(o.a)({},a.getResizerProps()),{},{className:[u.a.ResizeHandle,a.isResizing&&u.a.ResizeHandleActive].filter((function(e){return e})).join(" "),children:"\u22ee"}))]}))}))}))}))}),Object(h.jsx)("div",Object(o.a)(Object(o.a)({className:u.a.TableBody},d()),{},{children:Object(h.jsx)(g.a,{outerElementType:function(e,t){return Object(h.jsx)("div",Object(o.a)(Object(o.a)({},e),{},{style:Object(o.a)(Object(o.a)({},e.style),{},{overflowY:"scroll"}),forwardedRef:t}))},itemCount:j.length,height:H-56,itemSize:34,children:function(e){var t=e.index,a=e.style,n=j[t];return O(n),Object(h.jsx)("div",Object(o.a)(Object(o.a)({className:u.a.Row},n.getRowProps({style:a})),{},{children:n.cells.map((function(e){return Object(h.jsx)("div",Object(o.a)(Object(o.a)({className:u.a.BodyCell},e.getCellProps()),{},{children:e.render("Cell")}))}))}))}})}))]}))})]})}a(26);var C=function(){var e=Object(s.a)();return Object(h.jsx)("main",{className:"App",children:Object(h.jsx)(x,{data:e})})},_=function(e){e&&e instanceof Function&&a.e(3).then(a.bind(null,28)).then((function(t){var a=t.getCLS,n=t.getFID,c=t.getFCP,l=t.getLCP,r=t.getTTFB;a(e),n(e),c(e),l(e),r(e)}))};r.a.render(Object(h.jsx)(c.a.StrictMode,{children:Object(h.jsx)(C,{})}),document.getElementById("root")),_()}],[[27,1,2]]]);
//# sourceMappingURL=main.bc78bca4.chunk.js.map