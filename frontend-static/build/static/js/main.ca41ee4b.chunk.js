(this["webpackJsonpwarehouse-client"]=this["webpackJsonpwarehouse-client"]||[]).push([[0],{143:function(t,e,a){"use strict";a.r(e);var n=a(8),r=(a(67),a(68),a(0)),c=a(14),i=a.n(c),o=a(31),s=a(13),u=a.n(s),l=a(17),d=a(18),b=a.n(d),f="/products/gloves/",p={getAll:function(){var t=Object(l.a)(u.a.mark((function t(){var e;return u.a.wrap((function(t){for(;;)switch(t.prev=t.next){case 0:return console.log(f),t.next=3,b.a.get(f);case 3:return null===(e=t.sent).data&&(e.data=[]),t.abrupt("return",e.data);case 6:case"end":return t.stop()}}),t)})));return function(){return t.apply(this,arguments)}}()},j="/products/facemasks/",h={getAll:function(){var t=Object(l.a)(u.a.mark((function t(){var e;return u.a.wrap((function(t){for(;;)switch(t.prev=t.next){case 0:return console.log(j),t.next=3,b.a.get(j);case 3:return null===(e=t.sent).data&&(e.data=[]),t.abrupt("return",e.data);case 6:case"end":return t.stop()}}),t)})));return function(){return t.apply(this,arguments)}}()},v="/products/beanies/",g={getAll:function(){var t=Object(l.a)(u.a.mark((function t(){var e;return u.a.wrap((function(t){for(;;)switch(t.prev=t.next){case 0:return console.log(v),t.next=3,b.a.get(v);case 3:return null===(e=t.sent).data&&(e.data=[]),t.abrupt("return",e.data);case 6:case"end":return t.stop()}}),t)})));return function(){return t.apply(this,arguments)}}()},O=a(62),y=a(63),x=a(30),m=a(28),k=a.n(m),w=a(29),F=a.n(w),A=function(){var t=Object(r.useState)([]),e=Object(o.a)(t,2),a=e[0],c=e[1],i=Object(r.useState)([]),s=Object(o.a)(i,2),u=s[0],l=s[1],d=Object(r.useState)([]),b=Object(o.a)(d,2),f=b[0],j=b[1];Object(r.useEffect)((function(){p.getAll().then((function(t){return c(t)}))}),[]),Object(r.useEffect)((function(){g.getAll().then((function(t){return l(t)}))}),[]),Object(r.useEffect)((function(){h.getAll().then((function(t){return j(t)}))}),[]);var v=[{dataField:"api_id",text:"ID"},{dataField:"name",text:"Name"},{dataField:"colors",text:"Colors",formatter:function(t){return t.join(", ")}},{dataField:"price",text:"Price"},{dataField:"manufacturer",text:"Manufacturer"},{dataField:"availability",text:"Availability",formatter:function(t){return""===t?"Please refresh":t},style:function(t,e,a,n){return"INSTOCK"===t?{backgroundColor:"#5cb85c"}:"LESSTHAN10"===t?{backgroundColor:"#f0ad4e"}:"OUTOFSTOCK"===t?{backgroundColor:"#d9534f"}:void 0}}];return Object(n.jsx)("div",{className:"App",children:Object(n.jsxs)(O.a,{fluid:!0,style:{maxWidth:"1300px"},children:[Object(n.jsx)("center",{children:Object(n.jsx)("h1",{children:"reaktor warehouse"})}),Object(n.jsxs)(y.a,{defaultActiveKey:"gloves",id:"uncontrolled-tab",children:[Object(n.jsx)(x.a,{eventKey:"gloves",title:"Gloves",children:Object(n.jsx)(k.a,{bootstrap4:!0,hover:!0,striped:!0,noDataIndication:function(){return"No data currently available. Try again shortly by pressing ctrl+R"},keyField:"id",data:a,columns:v,pagination:F()({})})}),Object(n.jsx)(x.a,{eventKey:"beanies",title:"Beanies",children:Object(n.jsx)(k.a,{bootstrap4:!0,hover:!0,striped:!0,noDataIndication:function(){return"No data currently available. Try again shortly by pressing ctrl+R"},keyField:"id",data:u,columns:v,pagination:F()({})})}),Object(n.jsx)(x.a,{eventKey:"facemasks",title:"Facemasks",children:Object(n.jsx)(k.a,{bootstrap4:!0,hover:!0,striped:!0,noDataIndication:function(){return"No data currently available. Try again shortly by pressing ctrl+R"},keyField:"id",data:f,columns:v,pagination:F()({})})})]})]})})};i.a.render(Object(n.jsx)(A,{}),document.getElementById("root"))}},[[143,1,2]]]);
//# sourceMappingURL=main.ca41ee4b.chunk.js.map