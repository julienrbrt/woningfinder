<template>
  <Hero>
    <div class="mt-6 sm:max-w-xl">
      <h1
        class="text-3xl font-extrabold text-wf-purple tracking-tight sm:text-4xl"
      >
        {{ title }}
      </h1>
      <p class="mt-6 text-lg text-gray-500">
        Hoi 👋 ik ben
        <a
          href="https://julien.rbrt.fr"
          target="_blank"
          class="underline hover:text-gray-900"
          >Julien Robert</a
        >. Als softwareontwikkelaar heb ik zelf lange tijd geworsteld met het
        vinden van een huurwoning. Elke dag opnieuw naar de websites van
        verschillende woningcorporaties en verhuurders gaan om huizen te zoeken,
        terwijl mijn inbox volstroomde met duizenden e-mails over nieuw
        woningaanbod... Het was enorm frustrerend!
      </p>

      <p class="mt-6 text-lg text-gray-500">
        Daarom besloot ik WoningFinder te maken, in eerste instantie alleen voor
        mezelf, zodat ik het reageren op woningen kon automatiseren.
      </p>

      <p class="mt-6 text-lg text-gray-500">
        Toen ik ontdekte dat veel andere mensen tegen dezelfde problemen
        aanliepen, besloot ik WoningFinder gratis beschikbaar te maken voor
        iedereen 😊 Het is ook open-source en je kunt de sourcecode vinden op
        <a
          href="https://github.com/julienrbrt/woningfinder"
          target="_blank"
          class="underline hover:text-gray-900"
          >GitHub</a
        >. Hopelijk maakt WoningFinder jouw zoektocht een stuk makkelijker!
      </p>
    </div>

    <!-- profile -->
    <template v-slot:illustration>
      <img
        src="julien.jpg"
        alt="Julien Robert"
        class="relative mx-auto max-h-72 rounded-lg"
      />
    </template>

    <div class="items-center inline-flex mt-5 space-x-4">
      <BackButton />
    </div>
  </Hero>
</template>

<script>
import { InformationCircleIcon } from '@vue-hero-icons/outline'

export default {
  components: {
    InformationCircleIcon,
  },
  data() {
    return {
      title: 'Over ons – ehm mij',
    }
  },
  head() {
    return {
      title: this.title,
    }
  },
}
</script>
