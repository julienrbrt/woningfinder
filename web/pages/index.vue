<template>
  <div class="bg-white">
    <Hero>
      <template v-slot:illustration>
        <IllustrationCitiesList />
      </template>

      <!-- registration alert -->
      <AlertOk
        class="mb-6"
        v-if="alert && $route.query.thanks"
        @click="hideAlert"
        alert="Gefeliciteerd, je zoekopdracht is ingesteld 🎉 je ontvangt een e-mail om je zoekopdracht te activeren met instructies over hoe je kunt beginnen met automatisch reageren. Krijg je geen mail binnen? Kijk dan ook even bij jouw ongewenste mail."
      >
      </AlertOk>

      <div class="mt-6 sm:max-w-xl">
        <h1
          class="text-4xl font-extrabold text-wf-purple tracking-tight sm:text-5xl"
        >
          Reageer <span class="text-green-700">automatisch</span> op
          huurwoningen
        </h1>
        <p class="mt-6 text-xl text-gray-500">
          Mis nooit meer een kans op jouw droomhuis. Vind je perfecte huurwoning
          zonder dagelijks alle woningwebsites zelf te bezoeken. WoningFinder
          reageert automatisch op woningen die passen bij jouw zoekopdracht.
        </p>
      </div>
      <div class="mt-5 sm:mt-8 sm:flex sm:justify-center lg:justify-start">
        <NuxtLink to="/start" class="btn btn-primary btn-block">
          Nu starten
        </NuxtLink>
      </div>
      <div class="sm:flex sm:justify-center lg:justify-start">
        <div class="min-w-0 text-sm text-gray-500 sm:py-3">
          Meer kansen, minder werk.
        </div>
      </div>
    </Hero>

    <!-- Features section -->
    <LandingFeature ref="feature" id="feature" />

    <!-- FAQ section -->
    <LandingFaq id="faq" />

    <!-- CTA section -->
    <LandingCta id="contact" />
  </div>
</template>

<script>
export default {
  data() {
    return {
      alert: true,
    }
  },
  methods: {
    hideAlert() {
      this.alert = false
    },
  },
}
</script>
