<template>
  <div v-if="modelValue" class="tutorial-overlay">
    <div class="tutorial-focus" :class="`tutorial-focus--${currentStep.key}`"></div>
    <section class="tutorial-card" :class="`tutorial-card--${currentStep.key}`">
      <p class="tutorial-step">Etape {{ currentIndex + 1 }} / {{ steps.length }}</p>
      <h2>{{ currentStep.title }}</h2>
      <p>{{ currentStep.body }}</p>
      <div class="tutorial-actions">
        <button class="btn-secondary" type="button" @click="close">Passer</button>
        <button v-if="currentIndex > 0" class="btn-secondary" type="button" @click="currentIndex -= 1">Retour</button>
        <button class="btn-main-action" type="button" @click="next">
          {{ currentIndex === steps.length - 1 ? "Terminer" : "Suivant" }}
        </button>
      </div>
    </section>
  </div>
</template>

<script setup>
import { computed, ref, watch } from "vue";

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  }
});

const emit = defineEmits(["update:modelValue", "complete"]);

const currentIndex = ref(0);
const steps = [
  {
    key: "score",
    title: "Score et impact",
    body: "Cette zone resume vos points, vos objets recycles et l'impact CO2 calcule depuis le back."
  },
  {
    key: "annonces",
    title: "Vos annonces",
    body: "Ici vous suivez vos annonces, leur validation, leur statut et les actions possibles."
  },
  {
    key: "planning",
    title: "Calendrier",
    body: "Le planning regroupe vos formations, evenements et rendez-vous. Cliquez sur une entree pour ouvrir son detail."
  },
  {
    key: "nav",
    title: "Navigation",
    body: "La navbar donne acces au panier, a la messagerie, au catalogue, aux services et a votre compte."
  }
];

const currentStep = computed(() => steps[currentIndex.value]);

watch(
  () => props.modelValue,
  (isOpen) => {
    if (isOpen) currentIndex.value = 0;
  }
);

function close() {
  emit("update:modelValue", false);
  emit("complete");
}

function next() {
  if (currentIndex.value < steps.length - 1) {
    currentIndex.value += 1;
    return;
  }
  close();
}
</script>

<style scoped>
.tutorial-overlay {
  position: fixed;
  inset: 0;
  z-index: 120;
  background: rgba(10, 18, 13, 0.58);
  backdrop-filter: blur(2px);
}

.tutorial-focus {
  position: fixed;
  border: 3px solid #7ed999;
  border-radius: 16px;
  box-shadow: 0 0 0 9999px rgba(10, 18, 13, 0.42), 0 0 32px rgba(126, 217, 153, 0.35);
  pointer-events: none;
}

.tutorial-focus--score {
  top: 180px;
  left: 36px;
  width: min(54vw, 720px);
  height: 190px;
}

.tutorial-focus--annonces {
  top: 390px;
  left: 36px;
  width: min(86vw, 1180px);
  height: 250px;
}

.tutorial-focus--planning {
  bottom: 120px;
  left: 36px;
  width: min(86vw, 1180px);
  height: 190px;
}

.tutorial-focus--nav {
  top: 10px;
  left: 18px;
  right: 18px;
  height: 92px;
}

.tutorial-card {
  position: fixed;
  width: min(430px, calc(100vw - 36px));
  border: 1px solid rgba(126, 217, 153, 0.34);
  border-radius: 18px;
  background: #fff;
  padding: 22px;
  box-shadow: 0 30px 80px rgba(0, 0, 0, 0.24);
}

.tutorial-card--score,
.tutorial-card--annonces {
  top: 180px;
  right: 38px;
}

.tutorial-card--planning {
  bottom: 110px;
  right: 38px;
}

.tutorial-card--nav {
  top: 120px;
  right: 38px;
}

.tutorial-step {
  margin: 0 0 8px;
  color: var(--brand-green);
  font-size: 0.75rem;
  font-weight: 800;
  letter-spacing: 0.14em;
  text-transform: uppercase;
}

.tutorial-card h2 {
  margin: 0 0 8px;
  font-size: 1.35rem;
}

.tutorial-card p {
  color: var(--text-secondary);
  line-height: 1.55;
}

.tutorial-actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  margin-top: 18px;
}

@media (max-width: 760px) {
  .tutorial-focus {
    inset: 120px 16px auto;
    width: auto;
    height: 220px;
  }

  .tutorial-card {
    left: 18px;
    right: 18px;
    bottom: 18px;
    top: auto;
  }
}
</style>
