(window.webpackJsonp=window.webpackJsonp||[]).push([[41,6,7,23,24,25,26,27,28,29,30],{300:function(e,t,r){"use strict";r.r(t);var n=r(299),o={props:["alert"],components:{XCircleIcon:n.g,XIcon:n.h}},c=r(15),component=Object(c.a)(o,(function(){var e=this,t=e.$createElement,r=e._self._c||t;return r("div",{staticClass:"alert shadow-md alert-error"},[r("div",{staticClass:"text-gray-900"},[r("div",{staticClass:"w-max"},[r("XCircleIcon",{staticClass:"h-6 w-6"})],1),e._v(" "),r("span",{domProps:{innerHTML:e._s(e.alert)}})]),e._v(" "),r("div",{staticClass:"flex-none"},[r("button",{staticClass:"btn btn-sm btn-ghost",on:{click:function(t){return e.$emit("click")}}},[r("span",{staticClass:"sr-only"},[e._v("Ok")]),e._v(" "),r("XIcon")],1)])])}),[],!1,null,null,null);t.default=component.exports},301:function(e,t,r){"use strict";r.r(t);var n={props:["alert"]},o=r(15),component=Object(o.a)(n,(function(){var e=this,t=e.$createElement,r=e._self._c||t;return r("div",{staticClass:"alert mt-2 p-4"},[r("div",[r("div",{staticClass:"w-max"},[e._t("default")],2),e._v(" "),r("span",{staticClass:"text-sm font-medium text-gray-900",domProps:{innerHTML:e._s(e.alert)}})])])}),[],!1,null,null,null);t.default=component.exports},304:function(e,t,r){"use strict";r.r(t);var n=r(11),o=(r(67),r(53),r(68),r(23),r(54),r(141),r(299)),c=r(306)({accessToken:"pk.eyJ1Ijoid29uaW5nZmluZGVyIiwiYSI6ImNtM2YyNGRxczBpcHIyd3M4djZxeHVpbnAifQ.BpODv_QnbNXItkXb08FjtQ"}),l={components:{InformationCircleIcon:o.e,XIcon:o.h},props:["city"],data:function(){return{suggestedDistricts:this.city.district}},methods:{selectDistrict:function(input){var e=this;return Object(n.a)(regeneratorRuntime.mark((function t(){var r,i,n,o;return regeneratorRuntime.wrap((function(t){for(;;)switch(t.prev=t.next){case 0:if(r=[],!e.suggestedDistricts){t.next=8;break}if(0!=input.length){t.next=7;break}for(i=0;i<e.suggestedDistricts.length;i++)r.push(e.suggestedDistricts[i]);return t.abrupt("return",r);case 7:for(i=0;i<e.suggestedDistricts.length;i++)e.suggestedDistricts[i].toLowerCase().includes(input.toLowerCase())&&r.push(e.suggestedDistricts[i]);case 8:if(!(input.length>0)){t.next=14;break}return t.next=11,c.forwardGeocode({query:input,countries:["nl"],proximity:[e.city.longitude,e.city.latitude],types:["neighborhood","locality"],autocomplete:!0,language:["nl-NL"]}).send();case 11:for(n=t.sent,o=n.body,i=0;i<o.features.length;i++)o.features[i].place_name.includes(e.city.name)&&r.push(e.districtWithoutCity(o.features[i].place_name));case 14:return t.abrupt("return",r);case 15:case"end":return t.stop()}}),t)})))()},districtWithoutCity:function(e){return e.split(",",1)[0]},addCityDistrict:function(e){e&&(this.$store.commit("register/addCityDistrict",{city:this.city,district:e}),this.$refs.autocomplete.setValue(""),document.activeElement.blur(),this.$forceUpdate())},removeCityDistrict:function(e){e&&this.$store.commit("register/removeCityDistrict",{city:this.city,district:e}),this.$forceUpdate()}},computed:{getCity:function(){var e=this.$store.getters["register/getCity"](this.city);return e.district||(e.district=[]),e}}},m=r(15),component=Object(m.a)(l,(function(){var e=this,t=e.$createElement,r=e._self._c||t;return r("div",[e.city.district?r("p",{staticClass:"mt-2 text-base text-gray-500"},[e._v("\n    Er zijn "+e._s(e.city.district.length)+" voorgestelde wijken voor deze stad.\n  ")]):e._e(),e._v(" "),r("autocomplete",{ref:"autocomplete",staticClass:"mt-4",attrs:{search:e.selectDistrict,type:"text",placeholder:"Zoek buurt, wijk, etc.","aria-label":"Zoek buurt, wijk, etc.","debounce-time":200,"auto-select":""},on:{submit:e.addCityDistrict}}),e._v(" "),e.getCity.district.length>0?r("div",{staticClass:"mt-6 space-y-4"},e._l(e.getCity.district,(function(t){return r("div",{key:t,staticClass:"\n        relative\n        flex\n        items-center\n        rounded-lg\n        border border-gray-400\n        bg-white\n        shadow-sm\n        px-6\n        py-2\n        justify-between\n      "},[r("p",{staticClass:"text-sm font-medium text-gray-900"},[e._v("\n        "+e._s(t)+"\n      ")]),e._v(" "),r("button",{staticClass:"\n          inline-flex\n          rounded-md\n          p-1.5\n          text-gray-300\n          hover:text-red-300\n          focus:outline-none\n        ",attrs:{type:"button"},on:{click:function(r){return e.removeCityDistrict(t)}}},[r("XIcon",{attrs:{size:"1.5x"}})],1)])})),0):r("AlertInfo",{staticClass:"mt-4",attrs:{alert:"Je hebt nog geen wijk voorkeur. WoningFinder reageert daarom over de hele stad."}},[r("InformationCircleIcon",{staticClass:"h-5 w-5 text-gray-400"})],1)],1)}),[],!1,null,null,null);t.default=component.exports;installComponents(component,{AlertInfo:r(301).default})},309:function(e,t,r){"use strict";r.r(t);r(326),r(17),r(31),r(327),r(328),r(329),r(330),r(331),r(332),r(333),r(334),r(335),r(336),r(337),r(338),r(339),r(32),r(55),r(340),r(23),r(36),r(142);var n=r(299),o={components:{XIcon:n.h,InformationCircleIcon:n.e},props:["supported_cities"],data:function(){return{selectedCities:[],citiesList:new Map(this.supported_cities.sort((function(a,b){return a.name>b.name?1:-1})).map((function(e){return[e.name,e]}))),error:!1}},methods:{selectCity:function(input){var e=[];return 0==input.length?this.citiesList.forEach((function(t,r){e.push(r)})):this.citiesList.forEach((function(t,r){r.toLowerCase().startsWith(input.toLowerCase())&&e.push(r)})),e},addCity:function(e){e&&(this.$store.commit("register/addCity",e),this.$refs.autocomplete.setValue(""),document.activeElement.blur())},removeCity:function(e){this.$store.commit("register/removeCity",e)},validate:function(){return 0!=this.$store.getters["register/getCities"].length||(this.error=!0,!1)}},computed:{citiesSelection:function(){return this.$store.getters["register/getCities"]},hasSelection:function(){return this.citiesSelection.length>0}}},c=o,l=r(15),component=Object(l.a)(c,(function(){var e=this,t=e.$createElement,r=e._self._c||t;return r("div",{staticClass:"sm:max-w-xl"},[r("p",{staticClass:"mt-6 text-lg text-gray-500"},[e._v("Waar zoek je jouw woning?")]),e._v(" "),r("p",{staticClass:"mt-2 text-base text-gray-500"},[e._v("\n    Je kunt je gewenste steden selecteren en invullen tussen ons\n    "+e._s(e.supported_cities.length)+" beschikbare steden.\n  ")]),e._v(" "),r("autocomplete",{ref:"autocomplete",staticClass:"mt-4",attrs:{search:e.selectCity,type:"text",placeholder:"Steden selecteren of invullen","aria-label":"Steden selecteren of invullen","debounce-time":200,"auto-select":""},on:{submit:e.addCity}}),e._v(" "),e.hasSelection?e._e():r("AlertInfo",{staticClass:"mt-4",class:e.error?"bg-red-50":"bg-gray-50",attrs:{alert:"Je hebt geen steden geselecteerd."}},[r("InformationCircleIcon",{staticClass:"h-5 w-5",class:e.error?"text-red-400":"text-gray-400"})],1),e._v(" "),r("div",{staticClass:"mt-4 space-y-4"},e._l(e.citiesSelection,(function(t){return r("div",{key:t.name,staticClass:"\n        relative\n        flex\n        items-center\n        rounded-lg\n        border border-gray-400\n        bg-white\n        shadow-sm\n        px-6\n        py-2\n        justify-between\n      "},[r("p",{staticClass:"text-sm font-medium text-gray-900"},[e._v("\n        "+e._s(t.name)+"\n      ")]),e._v(" "),r("button",{staticClass:"\n          inline-flex\n          rounded-md\n          p-1.5\n          text-gray-300\n          hover:text-red-300\n          focus:outline-none\n        ",attrs:{type:"button"},on:{click:function(r){return e.removeCity(t)}}},[r("XIcon",{attrs:{size:"1.5x"}})],1)])})),0),e._v(" "),r("AlertInfo",{staticClass:"mt-4",attrs:{alert:'Staat je stad er niet tussen? Schrijf je in op onze\n             <a\n              to="/wachtlijst"\n              class="underline text-gray-700 hover:text-gray-900"\n              >wachtlijst</a>\n            en we laten je weten wanneer we jouw stad toegevoegd hebben.'}},[r("InformationCircleIcon",{staticClass:"h-5 w-5 text-gray-400"})],1)],1)}),[],!1,null,null,null);t.default=component.exports;installComponents(component,{AlertInfo:r(301).default})},310:function(e,t,r){"use strict";r.r(t);var n={props:["selected_cities"]},o=r(15),component=Object(o.a)(n,(function(){var e=this,t=e.$createElement,r=e._self._c||t;return r("div",{staticClass:"sm:max-w-xl"},[r("p",{staticClass:"mt-6 text-lg text-gray-500"},[e._v("\n    Je kan WoningFinder automatisch laten reageren op woningen in de wijk\n    en/of buurt die je wilt.\n  ")]),e._v(" "),r("p",{staticClass:"mt-2 text-base text-gray-500"},[e._v("\n    Als je een wijk voorkeur hebt, kun je dat nu zoeken, anders kun je deze\n    stap overslaan.\n  ")]),e._v(" "),r("div",{staticClass:"mt-6 space-y-4"},e._l(e.selected_cities,(function(t){return r("div",{key:t.name},[r("h2",{staticClass:"text-base font-medium text-gray-900"},[e._v("\n        "+e._s(t.name)+"\n      ")]),e._v(" "),r("RegisterCityDistrictPicker",{attrs:{city:t}})],1)})),0)])}),[],!1,null,null,null);t.default=component.exports;installComponents(component,{RegisterCityDistrictPicker:r(304).default})},311:function(e,t,r){"use strict";r.r(t);r(43);var n={props:["supported_housing"],data:function(){return{error:!1}},methods:{validate:function(){return 0!=this.$store.getters["register/getHousingType"].length||(this.error=!0,!1)},hideAlert:function(){this.error=!1},housingTitle:function(e){return"house"==e&&(e="huis"),e.charAt(0).toUpperCase()+e.slice(1)}},computed:{selectedType:{get:function(){return this.$store.getters["register/getHousingType"]},set:function(e){this.$store.commit("register/setHousingType",e)}}}},o=r(15),component=Object(o.a)(n,(function(){var e=this,t=e.$createElement,r=e._self._c||t;return r("div",{staticClass:"sm:max-w-xl"},[r("p",{staticClass:"mt-6 text-lg text-gray-500"},[e._v("Welke woningtype zoek je?")]),e._v(" "),r("p",{staticClass:"mt-2 text-base text-gray-500"},[e._v("\n    We weten nu genoeg over jouw gewenste locatie. We hebben nu een paar\n    vragen om meer te weten over de woning die jij zoekt.\n  ")]),e._v(" "),e.error?r("AlertError",{staticClass:"mt-4",attrs:{alert:"Selecteer een woningtype."},on:{click:e.hideAlert}}):e._e(),e._v(" "),r("fieldset",[r("legend",{staticClass:"sr-only"},[e._v("Woningtype")]),e._v(" "),r("div",{staticClass:"mt-6 space-y-4"},e._l(e.supported_housing,(function(t){return r("label",{key:t,staticClass:"\n          relative\n          block\n          rounded-lg\n          border\n          bg-white\n          shadow-sm\n          px-6\n          py-4\n          cursor-pointer\n          md:hover:border-wf-orange\n          sm:flex\n          sm:justify-between\n          focus-within:ring-1 focus-within:ring-wf-orange\n        ",class:[e.selectedType.indexOf(t)>-1?"border-wf-orange ":"border-gray-300"]},[r("input",{directives:[{name:"model",rawName:"v-model",value:e.selectedType,expression:"selectedType"}],staticClass:"sr-only",attrs:{type:"checkbox",name:"type"},domProps:{value:t,checked:Array.isArray(e.selectedType)?e._i(e.selectedType,t)>-1:e.selectedType},on:{change:function(r){var n=e.selectedType,o=r.target,c=!!o.checked;if(Array.isArray(n)){var l=t,m=e._i(n,l);o.checked?m<0&&(e.selectedType=n.concat([l])):m>-1&&(e.selectedType=n.slice(0,m).concat(n.slice(m+1)))}else e.selectedType=c}}}),e._v(" "),r("div",{staticClass:"flex items-center"},[r("div",{staticClass:"text-sm"},[r("p",{staticClass:"font-medium text-gray-900"},[e._v("\n              "+e._s(e.housingTitle(t))+"\n            ")])])])])})),0)])],1)}),[],!1,null,null,null);t.default=component.exports;installComponents(component,{AlertError:r(300).default})},312:function(e,t,r){"use strict";r.r(t);var n=r(299),o={components:{InformationCircleIcon:n.e,CurrencyEuroIcon:n.c},data:function(){return{error:!1,errorMsg:"",preferences:this.$store.getters["register/getHousingPreferencesDetails"]}},methods:{validate:function(){return this.preferences.maximum_price<=0?(this.error=!0,this.errorMsg="De maximale huurprijs moet hoger dan €0,- zijn.",!1):(this.preferences.number_bedroom=parseInt(this.preferences.number_bedroom),this.preferences.maximum_price=parseFloat(this.preferences.maximum_price),isNaN(this.preferences.number_bedroom)||isNaN(this.preferences.maximum_price)?(this.error=!0,this.errorMsg="Het aantal slaapkamers en/of de maximale huurprijs moeten een getal zijn.",!1):(this.$store.commit("register/setHousingPreferencesDetails",this.preferences),!0))},hideAlert:function(){this.error=!1}}},c=r(15),component=Object(c.a)(o,(function(){var e=this,t=e.$createElement,r=e._self._c||t;return r("div",{staticClass:"sm:max-w-xl"},[r("p",{staticClass:"mt-6 text-lg text-gray-500"},[e._v("\n    Vertel ons meer over jouw droomhuis\n  ")]),e._v(" "),r("p",{staticClass:"mt-2 text-base text-gray-500"},[e._v("\n    WoningFinder reageert niet op woningen die niet aan jouw voorkeuren\n    voldoen. Denk dus goed na.\n  ")]),e._v(" "),e.error?r("AlertError",{staticClass:"mt-4",attrs:{alert:e.errorMsg},on:{click:e.hideAlert}}):e._e(),e._v(" "),r("div",{staticClass:"mt-6 space-y-4"},[r("h2",{staticClass:"text-base font-medium text-gray-900"},[e._v("Binnen")]),e._v(" "),r("label",{staticClass:"block text-sm font-medium text-gray-900",attrs:{for:"bedroom"}},[e._v("Aantal slaapkamers")]),e._v(" "),r("select",{directives:[{name:"model",rawName:"v-model",value:e.preferences.number_bedroom,expression:"preferences.number_bedroom"}],staticClass:"\n        mt-2\n        shadow-sm\n        focus:ring-wf-orange focus:border-wf-orange\n        block\n        w-full\n        sm:text-sm\n        border-gray-300\n        rounded-md\n      ",attrs:{name:"bedroom"},on:{change:function(t){var r=Array.prototype.filter.call(t.target.options,(function(e){return e.selected})).map((function(e){return"_value"in e?e._value:e.value}));e.$set(e.preferences,"number_bedroom",t.target.multiple?r:r[0])}}},e._l([0,1,2,3,4,5],(function(t){return r("option",{key:t,domProps:{value:t}},[e._v("\n        "+e._s(t)+"+\n      ")])})),0),e._v(" "),r("div",{staticClass:"relative flex items-start"},[r("div",{staticClass:"flex items-center h-5"},[r("input",{directives:[{name:"model",rawName:"v-model",value:e.preferences.is_accessible,expression:"preferences.is_accessible"}],staticClass:"\n            focus:ring-wf-orange\n            h-4\n            w-4\n            text-wf-orange-dark\n            border-gray-300\n            rounded\n          ",attrs:{id:"accessible",name:"accessible",type:"checkbox"},domProps:{checked:Array.isArray(e.preferences.is_accessible)?e._i(e.preferences.is_accessible,null)>-1:e.preferences.is_accessible},on:{change:function(t){var r=e.preferences.is_accessible,n=t.target,o=!!n.checked;if(Array.isArray(r)){var c=e._i(r,null);n.checked?c<0&&e.$set(e.preferences,"is_accessible",r.concat([null])):c>-1&&e.$set(e.preferences,"is_accessible",r.slice(0,c).concat(r.slice(c+1)))}else e.$set(e.preferences,"is_accessible",o)}}})]),e._v(" "),r("label",{staticClass:"block text-sm font-medium text-gray-900 ml-3",attrs:{for:"accessible"}},[e._v("Woning toegankelijk (rolstoel of rolator)")])]),e._v(" "),r("h2",{staticClass:"text-base font-medium text-gray-900"},[e._v("Buiten")]),e._v(" "),r("div",{staticClass:"grid grid-flow-col grid-cols-2 grid-rows-2 gap-4"},[r("div",{staticClass:"relative flex items-start"},[r("div",{staticClass:"flex items-center h-5"},[r("input",{directives:[{name:"model",rawName:"v-model",value:e.preferences.has_balcony,expression:"preferences.has_balcony"}],staticClass:"\n              focus:ring-wf-orange\n              h-4\n              w-4\n              text-wf-orange-dark\n              border-gray-300\n              rounded\n            ",attrs:{id:"balcony",name:"balcony",type:"checkbox"},domProps:{checked:Array.isArray(e.preferences.has_balcony)?e._i(e.preferences.has_balcony,null)>-1:e.preferences.has_balcony},on:{change:function(t){var r=e.preferences.has_balcony,n=t.target,o=!!n.checked;if(Array.isArray(r)){var c=e._i(r,null);n.checked?c<0&&e.$set(e.preferences,"has_balcony",r.concat([null])):c>-1&&e.$set(e.preferences,"has_balcony",r.slice(0,c).concat(r.slice(c+1)))}else e.$set(e.preferences,"has_balcony",o)}}})]),e._v(" "),r("label",{staticClass:"block text-sm font-medium text-gray-900 ml-3",attrs:{for:"balcony"}},[e._v("Balkon")])]),e._v(" "),r("div",{staticClass:"relative flex items-start"},[r("div",{staticClass:"flex items-center h-5"},[r("input",{directives:[{name:"model",rawName:"v-model",value:e.preferences.has_garden,expression:"preferences.has_garden"}],staticClass:"\n              focus:ring-wf-orange\n              h-4\n              w-4\n              text-wf-orange-dark\n              border-gray-300\n              rounded\n            ",attrs:{id:"garden",name:"garden",type:"checkbox"},domProps:{checked:Array.isArray(e.preferences.has_garden)?e._i(e.preferences.has_garden,null)>-1:e.preferences.has_garden},on:{change:function(t){var r=e.preferences.has_garden,n=t.target,o=!!n.checked;if(Array.isArray(r)){var c=e._i(r,null);n.checked?c<0&&e.$set(e.preferences,"has_garden",r.concat([null])):c>-1&&e.$set(e.preferences,"has_garden",r.slice(0,c).concat(r.slice(c+1)))}else e.$set(e.preferences,"has_garden",o)}}})]),e._v(" "),r("label",{staticClass:"block text-sm font-medium text-gray-900 ml-3",attrs:{for:"garden"}},[e._v("Tuin")])]),e._v(" "),r("div",{staticClass:"relative flex items-start"},[r("div",{staticClass:"flex items-center h-5"},[r("input",{directives:[{name:"model",rawName:"v-model",value:e.preferences.has_garage,expression:"preferences.has_garage"}],staticClass:"\n              focus:ring-wf-orange\n              h-4\n              w-4\n              text-wf-orange-dark\n              border-gray-300\n              rounded\n            ",attrs:{id:"garage",name:"garage",type:"checkbox"},domProps:{checked:Array.isArray(e.preferences.has_garage)?e._i(e.preferences.has_garage,null)>-1:e.preferences.has_garage},on:{change:function(t){var r=e.preferences.has_garage,n=t.target,o=!!n.checked;if(Array.isArray(r)){var c=e._i(r,null);n.checked?c<0&&e.$set(e.preferences,"has_garage",r.concat([null])):c>-1&&e.$set(e.preferences,"has_garage",r.slice(0,c).concat(r.slice(c+1)))}else e.$set(e.preferences,"has_garage",o)}}})]),e._v(" "),r("label",{staticClass:"block text-sm font-medium text-gray-900 ml-3",attrs:{for:"garage"}},[e._v("Garage")])]),e._v(" "),r("div",{staticClass:"relative flex items-start"},[r("div",{staticClass:"flex items-center h-5"},[r("input",{directives:[{name:"model",rawName:"v-model",value:e.preferences.has_elevator,expression:"preferences.has_elevator"}],staticClass:"\n              focus:ring-wf-orange\n              h-4\n              w-4\n              text-wf-orange-dark\n              border-gray-300\n              rounded\n            ",attrs:{id:"elevator",name:"elevator",type:"checkbox"},domProps:{checked:Array.isArray(e.preferences.has_elevator)?e._i(e.preferences.has_elevator,null)>-1:e.preferences.has_elevator},on:{change:function(t){var r=e.preferences.has_elevator,n=t.target,o=!!n.checked;if(Array.isArray(r)){var c=e._i(r,null);n.checked?c<0&&e.$set(e.preferences,"has_elevator",r.concat([null])):c>-1&&e.$set(e.preferences,"has_elevator",r.slice(0,c).concat(r.slice(c+1)))}else e.$set(e.preferences,"has_elevator",o)}}})]),e._v(" "),r("label",{staticClass:"block text-sm font-medium text-gray-900 ml-3",attrs:{for:"elevator"}},[e._v("Lift")])])]),e._v(" "),r("h2",{staticClass:"text-base font-medium text-gray-900"},[e._v("Geld")]),e._v(" "),r("div",[r("label",{staticClass:"block text-sm font-medium text-gray-900",attrs:{for:"price"}},[e._v("Maximaal huurprijs")]),e._v(" "),r("div",{staticClass:"relative mt-4 rounded-md shadow-sm"},[e._m(0),e._v(" "),r("input",{directives:[{name:"model",rawName:"v-model",value:e.preferences.maximum_price,expression:"preferences.maximum_price"}],staticClass:"\n            focus:ring-wf-orange focus:border-wf-orange\n            block\n            w-full\n            pl-7\n            pr-12\n            sm:text-sm\n            border-gray-300\n            rounded-md\n          ",attrs:{type:"number",name:"price",id:"price",placeholder:"1000"},domProps:{value:e.preferences.maximum_price},on:{input:function(t){t.target.composing||e.$set(e.preferences,"maximum_price",t.target.value)}}}),e._v(" "),e._m(1)])])]),e._v(" "),r("AlertInfo",{staticClass:"mt-4",attrs:{alert:'Heb jij een wens die niet in ons zoekopdracht staat? Neem dan\n          <a\n            to="/contact"\n            class="underline text-gray-700 hover:text-gray-900"\n            >contact</a>\n          met ons op.'}},[r("InformationCircleIcon",{staticClass:"h-5 w-5 text-gray-400"})],1)],1)}),[function(){var e=this,t=e.$createElement,r=e._self._c||t;return r("div",{staticClass:"\n            absolute\n            inset-y-0\n            left-0\n            pl-3\n            flex\n            items-center\n            pointer-events-none\n          "},[r("span",{staticClass:"text-gray-500 sm:text-sm"},[e._v(" € ")])])},function(){var e=this,t=e.$createElement,r=e._self._c||t;return r("div",{staticClass:"\n            absolute\n            inset-y-0\n            right-0\n            pr-3\n            flex\n            items-center\n            pointer-events-none\n          "},[r("span",{staticClass:"text-gray-500 sm:text-sm"},[e._v(" per maand ")])])}],!1,null,null,null);t.default=component.exports;installComponents(component,{AlertError:r(300).default,AlertInfo:r(301).default})},313:function(e,t,r){"use strict";r.r(t);var n={props:["buttonStr","currentStep","totalStep"]},o=r(15),component=Object(o.a)(n,(function(){var e=this,t=e.$createElement,r=e._self._c||t;return r("div",{staticClass:"items-center w-full inline-flex mt-5"},[r("div",{directives:[{name:"show",rawName:"v-show",value:e.currentStep<e.totalStep,expression:"currentStep < totalStep"}],staticClass:"mr-4"},[1==e.currentStep?r("div",[r("BackButton")],1):r("div",[r("a",{staticClass:"btn btn-ghost text-gray-500",on:{click:function(t){return e.$emit("previous")}}},[e._v("\n        Terug\n      ")])])]),e._v(" "),r("a",{staticClass:"btn btn-primary",class:e.currentStep==e.totalStep?"flex-1":"max-w-min",on:{click:function(t){return e.$emit("validate")}}},[e._v("\n    "+e._s(e.buttonStr)+"\n  ")]),e._v(" "),r("p",{staticClass:"\n      flex-1\n      whitespace-nowrap\n      text-sm\n      font-medium\n      text-gray-500 text-right\n    "},[e._v("\n    "+e._s(e.currentStep)+" / "+e._s(e.totalStep)+"\n  ")])])}),[],!1,null,null,null);t.default=component.exports;installComponents(component,{BackButton:r(202).default})},357:function(e,t,r){"use strict";r.r(t);r(23);var n={components:{InformationCircleIcon:r(299).e},data:function(){return{error:!1,errorMsg:"",customer:this.$store.getters["register/getCustomer"]}},methods:{validate:function(){return this.customer.yearly_income=parseInt(this.customer.yearly_income),this.customer.family_size=parseInt(this.customer.family_size),this.customer.birth_year=parseInt(this.customer.birth_year),isNaN(this.customer.yearly_income)||isNaN(this.customer.family_size)||""==this.customer.name||""==this.customer.email||this.customer.birth_year<1920?(this.error=!0,this.errorMsg="We hebben al je gegevens nodig. Controleer het formulier nogmaals.",!1):(this.$store.commit("register/setCustomer",this.customer),!0)},hideAlert:function(){this.error=!1}},computed:{currentYear:function(){return(new Date).getFullYear()},showHasChildren:function(){return this.currentYear-this.customer.birth_year<23}}},o=r(15),component=Object(o.a)(n,(function(){var e=this,t=e.$createElement,r=e._self._c||t;return r("div",{staticClass:"sm:max-w-xl"},[r("p",{staticClass:"mt-6 text-lg text-gray-500"},[e._v("Laten we het persoonlijker maken")]),e._v(" "),e.error?r("AlertError",{staticClass:"mt-4",attrs:{alert:e.errorMsg},on:{click:e.hideAlert}}):e._e(),e._v(" "),r("div",{staticClass:"mt-6 space-y-4"},[r("div",{staticClass:"grid grid-cols-6 gap-6"},[r("div",{staticClass:"col-span-6 sm:col-span-3"},[r("label",{staticClass:"block text-sm font-medium text-gray-900",attrs:{for:"name"}},[e._v("\n          Naam\n        ")]),e._v(" "),r("input",{directives:[{name:"model",rawName:"v-model",value:e.customer.name,expression:"customer.name"}],staticClass:"\n            mt-2\n            shadow-sm\n            focus:ring-wf-orange focus:border-wf-orange\n            block\n            w-full\n            sm:text-sm\n            border-gray-300\n            rounded-md\n          ",attrs:{type:"text",name:"name",id:"name",autocomplete:"given-name"},domProps:{value:e.customer.name},on:{input:function(t){t.target.composing||e.$set(e.customer,"name",t.target.value)}}})]),e._v(" "),r("div",{staticClass:"col-span-6 sm:col-span-3"},[r("label",{staticClass:"block text-sm font-medium text-gray-900",attrs:{for:"birth_year"}},[e._v("\n          Geboorte jaar\n        ")]),e._v(" "),r("input",{directives:[{name:"model",rawName:"v-model",value:e.customer.birth_year,expression:"customer.birth_year"}],staticClass:"\n            mt-2\n            shadow-sm\n            focus:ring-wf-orange focus:border-wf-orange\n            block\n            w-full\n            sm:text-sm\n            border-gray-300\n            rounded-md\n          ",attrs:{name:"birth_year",type:"number",min:"1920",max:e.currentYear,step:"1"},domProps:{value:e.customer.birth_year},on:{input:function(t){t.target.composing||e.$set(e.customer,"birth_year",t.target.value)}}})])]),e._v(" "),r("label",{staticClass:"block text-sm font-medium text-gray-900",attrs:{for:"email"}},[e._v("\n      E-mailadres\n    ")]),e._v(" "),r("input",{directives:[{name:"model",rawName:"v-model",value:e.customer.email,expression:"customer.email"}],staticClass:"\n        shadow-sm\n        focus:ring-wf-orange focus:border-wf-orange\n        block\n        w-full\n        sm:text-sm\n        border-gray-300\n        rounded-md\n      ",attrs:{id:"email",name:"email",type:"email",autocomplete:"email"},domProps:{value:e.customer.email},on:{input:function(t){t.target.composing||e.$set(e.customer,"email",t.target.value)}}}),e._v(" "),r("div",{staticClass:"grid grid-cols-6 gap-6"},[r("div",{staticClass:"col-span-6 sm:col-span-3"},[r("label",{staticClass:"block text-sm font-medium text-gray-900",attrs:{for:"income"}},[e._v("Jaarlijks inkomen")]),e._v(" "),r("div",{staticClass:"mt-2 relative rounded-md shadow-sm"},[e._m(0),e._v(" "),r("input",{directives:[{name:"model",rawName:"v-model",value:e.customer.yearly_income,expression:"customer.yearly_income"}],staticClass:"\n              focus:ring-wf-orange focus:border-wf-orange\n              block\n              w-full\n              pl-7\n              pr-12\n              sm:text-sm\n              border-gray-300\n              rounded-md\n            ",attrs:{type:"number",name:"income",id:"income",min:"0",placeholder:"20000"},domProps:{value:e.customer.yearly_income},on:{input:function(t){t.target.composing||e.$set(e.customer,"yearly_income",t.target.value)}}}),e._v(" "),e._m(1)])]),e._v(" "),r("div",{staticClass:"col-span-6 sm:col-span-3"},[r("label",{staticClass:"block text-sm font-medium text-gray-900",attrs:{for:"family"}},[e._v("\n          Aantal inwoners (inclusief jij)\n        ")]),e._v(" "),r("input",{directives:[{name:"model",rawName:"v-model",value:e.customer.family_size,expression:"customer.family_size"}],staticClass:"\n            mt-2\n            shadow-sm\n            focus:ring-wf-orange focus:border-wf-orange\n            block\n            w-full\n            sm:text-sm\n            border-gray-300\n            rounded-md\n          ",attrs:{type:"number",name:"family",id:"family",min:"1",max:"12",placeholder:"1"},domProps:{value:e.customer.family_size},on:{input:function(t){t.target.composing||e.$set(e.customer,"family_size",t.target.value)}}})])]),e._v(" "),e.showHasChildren?r("div",[r("p",{staticClass:"text-sm text-gray-500"},[e._v("\n        Omdat je jonger bent dan 23, moet de maximaal huurprijs lager zijn om\n        huurtoeslag te krijgen. Als jij een kind hebt, geld dit niet voor jou.\n      ")]),e._v(" "),r("div",{staticClass:"mt-4 inline-flex"},[r("input",{directives:[{name:"model",rawName:"v-model",value:e.customer.has_children_same_housing,expression:"customer.has_children_same_housing"}],staticClass:"\n            focus:ring-wf-orange\n            h-4\n            w-4\n            text-wf-orange-dark\n            border-gray-300\n            rounded\n          ",attrs:{id:"children",name:"children",type:"checkbox"},domProps:{checked:Array.isArray(e.customer.has_children_same_housing)?e._i(e.customer.has_children_same_housing,null)>-1:e.customer.has_children_same_housing},on:{change:function(t){var r=e.customer.has_children_same_housing,n=t.target,o=!!n.checked;if(Array.isArray(r)){var c=e._i(r,null);n.checked?c<0&&e.$set(e.customer,"has_children_same_housing",r.concat([null])):c>-1&&e.$set(e.customer,"has_children_same_housing",r.slice(0,c).concat(r.slice(c+1)))}else e.$set(e.customer,"has_children_same_housing",o)}}}),e._v(" "),r("label",{staticClass:"block text-sm font-medium text-gray-900 ml-3",attrs:{for:"children"}},[e._v("Ik heb een meeverhuizende kindje")])])]):e._e(),e._v(" "),r("AlertInfo",{staticClass:"mt-4",attrs:{alert:'WoningFinder reageert alleen op\n            <a\n              class="underline text-gray-700 hover:text-gray-900"\n              target="_blank"\n              href="https://www.woningmarktbeleid.nl/onderwerpen/daeb/toewijzen-door-woningcorporaties/regels-voor-toewijzen"\n              >passende woningen</a\n            >. We hebben dus extra gegevens nodig zoals het aantal\n            medebewoners of je (gezamelijk) jaarlijks inkomen. Al je gegevens\n            blijven altijd privé en worden nooit gedeeld.'}},[r("InformationCircleIcon",{staticClass:"h-5 w-5 text-gray-400"})],1)],1)],1)}),[function(){var e=this,t=e.$createElement,r=e._self._c||t;return r("div",{staticClass:"\n              absolute\n              inset-y-0\n              left-0\n              pl-3\n              flex\n              items-center\n              pointer-events-none\n            "},[r("span",{staticClass:"text-gray-500 sm:text-sm"},[e._v(" € ")])])},function(){var e=this,t=e.$createElement,r=e._self._c||t;return r("div",{staticClass:"\n              absolute\n              inset-y-0\n              right-0\n              pr-3\n              flex\n              items-center\n              pointer-events-none\n            "},[r("span",{staticClass:"text-gray-500 sm:text-sm"},[e._v(" per jaar ")])])}],!1,null,null,null);t.default=component.exports;installComponents(component,{AlertError:r(300).default,AlertInfo:r(301).default})},358:function(e,t,r){"use strict";r.r(t);var n={components:{InformationCircleIcon:r(299).e}},o=r(15),component=Object(o.a)(n,(function(){var e=this,t=e.$createElement,r=e._self._c||t;return r("div",{staticClass:"sm:max-w-xl"},[r("p",{staticClass:"mt-6 text-lg text-gray-500"},[e._v("\n    Je zoekopdracht is bijna ingesteld, dit is wat jij moet weten\n  ")]),e._v(" "),r("AlertInfo",{staticClass:"mt-6",attrs:{alert:"Je reageert met WoningFinder op alle huurwoningen die aan je zoekopdracht voldoen."}},[r("InformationCircleIcon",{staticClass:"h-5 w-5 text-gray-400"})],1),e._v(" "),r("AlertInfo",{attrs:{alert:"WoningFinder kan alleen op beschikbare woningen reageren. Zijn er geen woningen\n      beschikbaar die matchen met jouw zoekopdracht dan kan WoningFinder ook nergens op\n      reageren. Je ontvangt wel altijd een wekelijkse update van ons."}},[r("InformationCircleIcon",{staticClass:"h-5 w-5 text-gray-400"})],1),e._v(" "),r("AlertInfo",{attrs:{alert:"WoningFinder reageert op de woningen die je zou willen (je\n      zoekopdracht). Maar uiteraard kunnen we niet garanderen dat jij sneller een huis\n      krijgt, je reageert wel automatisch, maar het zijn altijd de\n      verhuurders die de selecties doen (loting, wie eerst komt wie\n      eerst maalt of inschrijfperiode)."}},[r("InformationCircleIcon",{staticClass:"h-5 w-5 text-gray-400"})],1),e._v(" "),e._m(0)],1)}),[function(){var e=this,t=e.$createElement,r=e._self._c||t;return r("p",{staticClass:"mt-4 text-sm text-gray-900"},[e._v("\n    Door op Aanmelden te klikken, ga je akkoord met onze\n    "),r("a",{staticClass:"underline",attrs:{href:"/voorwaarden",target:"_blank"}},[e._v("Voorwaarden")]),e._v("\n    en\n    "),r("a",{staticClass:"underline",attrs:{href:"/voorwaarden/gebruiks",target:"_blank"}},[e._v("Gebruiksvoorwaarden")]),e._v(". Meer informatie over hoe we je gegevens gebruiken vind je in ons\n    "),r("a",{staticClass:"underline",attrs:{href:"/voorwaarden/privacy",target:"_blank"}},[e._v("Privacybeleid")]),e._v(".\n  ")])}],!1,null,null,null);t.default=component.exports;installComponents(component,{AlertInfo:r(301).default})},408:function(e,t,r){"use strict";r.r(t);var n=r(11),o=(r(36),r(204),r(23),r(67),{asyncData:function(e){return Object(n.a)(regeneratorRuntime.mark((function t(){var r,n;return regeneratorRuntime.wrap((function(t){for(;;)switch(t.prev=t.next){case 0:return r=e.$axios,t.next=3,r.$get("offering",{progress:!0});case 3:return n=t.sent,t.abrupt("return",{offering:n});case 5:case"end":return t.stop()}}),t)})))()},data:function(){return{title:"Je WoningFinder zoekopdracht",submitted:!1,error:!1,errorMsg:"Er is iets misgegaan. Controleer het formulier nogmaals. Blijf dit gebeuren? Neem dan contact met ons op.",currentStep:1,totalStep:6,selected_cities:[]}},head:function(){return{title:this.title}},methods:{buttonStr:function(){return this.currentStep==this.totalStep?"Aanmelden":"Volgende"},next:function(){this.currentStep<=this.totalStep&&this.currentStep++},previous:function(){this.currentStep>1&&this.currentStep--},validate:function(){var e=this;return Object(n.a)(regeneratorRuntime.mark((function t(){var data;return regeneratorRuntime.wrap((function(t){for(;;)switch(t.prev=t.next){case 0:t.t0=e.currentStep,t.next=1===t.t0?3:2===t.t0?6:3===t.t0?7:4===t.t0?10:5===t.t0?13:6===t.t0?16:27;break;case 3:if(e.$refs.registerCity.validate()){t.next=5;break}return t.abrupt("return");case 5:case 6:return t.abrupt("break",27);case 7:if(e.$refs.registerHousing.validate()){t.next=9;break}return t.abrupt("return");case 9:return t.abrupt("break",27);case 10:if(e.$refs.registerHousingPreferences.validate()){t.next=12;break}return t.abrupt("return");case 12:return t.abrupt("break",27);case 13:if(e.$refs.registerCustomer.validate()){t.next=15;break}return t.abrupt("return");case 15:return t.abrupt("break",27);case 16:if(e.submitted){t.next=26;break}return e.$nuxt.$loading.start(),data=e.$store.getters["register/getRegister"],t.next=21,e.submit(data);case 21:if(e.submitted=!0,!e.error){t.next=24;break}return t.abrupt("return");case 24:e.$nuxt.$loading.finish(),e.$router.push({path:"/",query:{thanks:!0}});case 26:return t.abrupt("return");case 27:e.next();case 28:case"end":return t.stop()}}),t)})))()},submit:function(input){var e=this;return Object(n.a)(regeneratorRuntime.mark((function t(){return regeneratorRuntime.wrap((function(t){for(;;)switch(t.prev=t.next){case 0:return t.next=2,e.$axios.$post("register",input).then((function(e){return e})).catch((function(t){e.error=!0,e.errorMsg='Er is iets misgegaan: "'+t.response.data.message+'".'}));case 2:case"end":return t.stop()}}),t)})))()},hideAlert:function(){this.submitted=!1,this.error=!1}},computed:{citiesSelection:function(){var e=this,t=this.$store.getters["register/getCities"],r=[];return t.forEach((function(t){r.push(e.offering.supported_cities.find((function(e){return e.name==t.name})))})),r}}}),c=r(15),component=Object(c.a)(o,(function(){var e=this,t=e.$createElement,r=e._self._c||t;return r("Hero",[r("div",{staticClass:"mt-6 sm:max-w-xl"},[r("h1",{staticClass:"text-3xl font-extrabold text-wf-purple tracking-tight"},[e._v("\n      "+e._s(e.title)+"\n    ")])]),e._v(" "),e.error?r("AlertError",{staticClass:"mt-4",attrs:{alert:e.errorMsg},on:{click:e.hideAlert}}):e._e(),e._v(" "),r("RegisterCity",{directives:[{name:"show",rawName:"v-show",value:1==e.currentStep,expression:"currentStep == 1"}],ref:"registerCity",attrs:{supported_cities:e.offering.supported_cities}}),e._v(" "),r("RegisterCityDistrict",{directives:[{name:"show",rawName:"v-show",value:2==e.currentStep,expression:"currentStep == 2"}],ref:"registerCityDistrict",attrs:{selected_cities:e.citiesSelection}}),e._v(" "),r("RegisterHousing",{directives:[{name:"show",rawName:"v-show",value:3==e.currentStep,expression:"currentStep == 3"}],ref:"registerHousing",attrs:{supported_housing:e.offering.supported_housing_types}}),e._v(" "),r("RegisterHousingPreferences",{directives:[{name:"show",rawName:"v-show",value:4==e.currentStep,expression:"currentStep == 4"}],ref:"registerHousingPreferences"}),e._v(" "),r("RegisterCustomer",{directives:[{name:"show",rawName:"v-show",value:5==e.currentStep,expression:"currentStep == 5"}],ref:"registerCustomer"}),e._v(" "),r("RegisterTerms",{directives:[{name:"show",rawName:"v-show",value:6==e.currentStep,expression:"currentStep == 6"}],ref:"registerTerms"}),e._v(" "),r("Navigator",{attrs:{currentStep:e.currentStep,totalStep:e.totalStep,buttonStr:e.buttonStr()},on:{validate:e.validate,previous:e.previous}})],1)}),[],!1,null,null,null);t.default=component.exports;installComponents(component,{AlertError:r(300).default,RegisterCity:r(309).default,RegisterCityDistrict:r(310).default,RegisterHousing:r(311).default,RegisterHousingPreferences:r(312).default,RegisterCustomer:r(357).default,RegisterTerms:r(358).default,Navigator:r(313).default,Hero:r(201).default})}}]);