(window.webpackJsonp=window.webpackJsonp||[]).push([[31,7,18],{301:function(t,n,e){"use strict";e.r(n);var o={props:["alert"]},r=e(15),component=Object(r.a)(o,(function(){var t=this,n=t.$createElement,e=t._self._c||n;return e("div",{staticClass:"alert mt-2 p-4"},[e("div",[e("div",{staticClass:"w-max"},[t._t("default")],2),t._v(" "),e("span",{staticClass:"text-sm font-medium text-gray-900",domProps:{innerHTML:t._s(t.alert)}})])])}),[],!1,null,null,null);n.default=component.exports},316:function(t,n,e){var map={"./dak-2.png":317,"./dak.png":318,"./dewoningzoeker.png":319,"./ikwilhuren.png":320,"./roomspot.png":321,"./woninghuren.png":322,"./woonnethaaglanden.png":323};function o(t){var n=r(t);return e(n)}function r(t){if(!e.o(map,t)){var n=new Error("Cannot find module '"+t+"'");throw n.code="MODULE_NOT_FOUND",n}return map[t]}o.keys=function(){return Object.keys(map)},o.resolve=r,t.exports=o,o.id=316},317:function(t,n,e){t.exports=e.p+"img/dak-2.fb822c2.png"},318:function(t,n,e){t.exports=e.p+"img/dak.3ceb542.png"},319:function(t,n,e){t.exports=e.p+"img/dewoningzoeker.cf75324.png"},320:function(t,n,e){t.exports=e.p+"img/ikwilhuren.b3d7ff2.png"},321:function(t,n,e){t.exports=e.p+"img/roomspot.80cc900.png"},322:function(t,n,e){t.exports=e.p+"img/woninghuren.8ee2cb1.png"},323:function(t,n,e){t.exports=e.p+"img/woonnethaaglanden.94f5613.png"},349:function(t,n,e){"use strict";e.r(n);var o={props:["name","image","website"]},r=e(15),component=Object(r.a)(o,(function(){var t=this,n=t.$createElement,o=t._self._c||n;return o("a",{staticClass:"cursor-pointer",attrs:{href:t.website,target:"_blank"}},[o("div",{staticClass:"\n      tooltip\n      col-span-1\n      flex\n      justify-center\n      py-8\n      px-8\n      bg-gray-50\n      rounded-lg\n    ",attrs:{"data-tip":t.name}},[o("img",{staticClass:"max-h-12",attrs:{src:e(316)("./"+t.image),alt:t.name}})])])}),[],!1,null,null,null);n.default=component.exports},403:function(t,n,e){"use strict";e.r(n);var o={components:{InformationCircleIcon:e(299).e},data:function(){return{title:"Beschikbare woningcorporaties en verhuurders"}},head:function(){return{title:this.title}}},r=e(15),component=Object(r.a)(o,(function(){var t=this,n=t.$createElement,e=t._self._c||n;return e("Hero",[e("div",{staticClass:"mt-6 sm:max-w-xl"},[e("h1",{staticClass:"text-3xl font-extrabold text-wf-purple tracking-tight sm:text-4xl"},[t._v("\n      "+t._s(t.title)+"\n    ")]),t._v(" "),e("p",{staticClass:"mt-6 text-lg text-gray-500"},[t._v("\n      Hier zijn de woningaanbod websites waarin je automatisch kan reageren\n      met WoningFinder.\n    ")]),t._v(" "),e("AlertInfo",{staticClass:"mt-6",attrs:{alert:'Zoek je een site die hier nog niet staat? Neem dan\n            <a\n              to="/contact"\n              class="underline text-gray-700 hover:text-gray-900"\n              >contact</a>\n            met ons op.'}},[e("InformationCircleIcon",{staticClass:"h-5 w-5 text-gray-400"})],1),t._v(" "),e("h2",{staticClass:"mt-4 mb-2 text-xl text-gray-900"},[t._v("Heel Nederland")]),t._v(" "),e("div",{staticClass:"mt-8 grid grid-cols-2 gap-0.5 md:grid-cols-3 lg:mt-0 lg:grid-cols-2"},[e("LandingCorporationLogo",{attrs:{name:"Ikwilhuren",image:"ikwilhuren.png",website:"https://ikwilhuren.nu"}}),t._v(" "),e("LandingCorporationLogo",{attrs:{name:"DĀK WoningNet",image:"dak.png",website:"https://www.woningnet.nl"}})],1),t._v(" "),e("h2",{staticClass:"mt-4 mb-2 text-xl text-gray-900"},[t._v("Regio Overijssel")]),t._v(" "),e("div",{staticClass:"mt-8 grid grid-cols-2 gap-0.5 md:grid-cols-3 lg:mt-0 lg:grid-cols-2"},[e("LandingCorporationLogo",{attrs:{name:"WoningHuren",image:"woninghuren.png",website:"https://www.woninghuren.nl"}}),t._v(" "),e("LandingCorporationLogo",{attrs:{name:"De WoningZoeker",image:"dewoningzoeker.png",website:"https://www.dewoningzoeker.nl"}}),t._v(" "),e("LandingCorporationLogo",{attrs:{name:"Roomspot",image:"roomspot.png",website:"https://www.roomspot.nl"}})],1),t._v(" "),e("h2",{staticClass:"mt-4 mb-2 text-xl text-gray-900"},[t._v("Regio Zuid-Holland")]),t._v(" "),e("div",{staticClass:"mt-8 grid grid-cols-2 gap-0.5 md:grid-cols-3 lg:mt-0 lg:grid-cols-2"},[e("LandingCorporationLogo",{attrs:{name:"Woonnet Haaglanden",image:"woonnethaaglanden.png",website:"https://www.woonnet-haaglanden.nl"}})],1),t._v(" "),e("AlertInfo",{staticClass:"mt-4",attrs:{alert:"Let op: WoningFinder is zelf geen woninganbood website maar een aggregator, alle namen en\n      logo's op deze pagina zijn eigendom van de verhuurders zelf."}},[e("InformationCircleIcon",{staticClass:"h-5 w-5 text-gray-400"})],1),t._v(" "),e("div",{staticClass:"mt-4 inline-flex flex-col sm:flex-row w-max items-start sm:items-center justify-center"},[e("p",{staticClass:"text-gray-500"},[t._v("Reageer je nog steeds niet automatisch?")]),t._v(" "),e("NuxtLink",{staticClass:"mt-2 sm:mt-0 sm:ml-4 w-auto py-2 btn btn-secondary",attrs:{to:"/start"}},[t._v("\n        Begin nu\n      ")])],1),t._v(" "),e("div",{staticClass:"items-center mt-5 space-x-4"},[e("BackButton")],1)],1)])}),[],!1,null,null,null);n.default=component.exports;installComponents(component,{AlertInfo:e(301).default,LandingCorporationLogo:e(349).default,BackButton:e(202).default,Hero:e(201).default})}}]);