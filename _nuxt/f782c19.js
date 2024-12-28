(window.webpackJsonp=window.webpackJsonp||[]).push([[13],{356:function(e,t,r){"use strict";r.r(t);r(144),r(23);var n={props:["customer","housing_preferences"],methods:{housingTypeTitle:function(){for(var e=[],i=0;i<this.housing_preferences.type.length;i++)switch(this.housing_preferences.type[i]){case"house":e.push("Huis");break;case"appartement":e.push("Appartement")}return e.join(" of ")},cityTitle:function(){for(var e=[],i=0;i<this.housing_preferences.city.length;i++)e.push(this.housing_preferences.city[i].name);return e.join(", ")},hasExtras:function(){return this.housing_preferences.is_accessible||this.housing_preferences.has_balcony||this.housing_preferences.has_garden||this.housing_preferences.has_elevator},translateTitle:function(){var e=[];return this.housing_preferences.is_accessible&&e.push("woning toegankelijkheid"),this.housing_preferences.has_balcony&&e.push("balkon"),this.housing_preferences.has_garden&&e.push("tuin"),this.housing_preferences.has_garage&&e.push("garage"),this.housing_preferences.has_elevator&&e.push("lift"),e.join(", ")}}},c=r(15),component=Object(c.a)(n,(function(){var e=this,t=e.$createElement,r=e._self._c||t;return r("div",{staticClass:"mt-6 card bg-white shadow"},[r("div",{staticClass:"card-body"},[r("h1",{staticClass:"card-title text-xl text-wf-puple"},[e._v("Je profiel")]),e._v(" "),r("ul",{staticClass:"text-gray-500 text-base grid grid-rows-2 grid-cols-2"},[r("li",{staticClass:"mb-2"},[e._v("👤 "+e._s(e.customer.name))]),e._v(" "),r("li",{staticClass:"mb-2 break-all"},[e._v("✉️ "+e._s(e.customer.email))]),e._v(" "),r("li",{staticClass:"mb-2"},[e._v("🗓 "+e._s(e.customer.birth_year))]),e._v(" "),r("li",{staticClass:"mb-2"},[e._v("💰 €"+e._s(e.customer.yearly_income)+" per jaar")]),e._v(" "),r("li",{staticClass:"mb-2"},[e._v("👪 "+e._s(e.customer.family_size)+" inwoner(s)")])]),e._v(" "),r("h1",{staticClass:"card-title text-xl text-wf-puple"},[e._v("\n      Je zoekopdracht in het kort\n    ")]),e._v(" "),r("ul",{staticClass:"text-gray-500 text-base grid grid-rows-2 grid-cols-2"},[r("li",{staticClass:"mb-2"},[e._v("🏠 "+e._s(e.housingTypeTitle()))]),e._v(" "),r("li",{staticClass:"mb-2"},[e._v("📍 "+e._s(e.cityTitle()))]),e._v(" "),e.housing_preferences.maximum_price>0?r("li",{staticClass:"mb-2"},[e._v("\n        💰 €"+e._s(e.housing_preferences.maximum_price)+" p/m maximum\n      ")]):e._e(),e._v(" "),r("li",{staticClass:"mb-2"},[e._v("\n        🛏 Minimaal\n        "+e._s(e.housing_preferences.number_bedroom)+" slaapkamer(s)\n      ")]),e._v(" "),e.hasExtras()?r("li",{staticClass:"mb-2"},[e._v("\n        🪴 Extras zoals "+e._s(e.translateTitle())+"\n      ")]):e._e()])])])}),[],!1,null,null,null);t.default=component.exports}}]);