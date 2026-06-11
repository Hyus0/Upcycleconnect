<template>
  <section class="stack">
    <div class="page-heading">
      <div>
        <div class="eyebrow">Pilotage complet</div>
        <h2 class="page-title">Back connecté au front</h2>
        <p class="page-subtitle">
          Vue centralisée des fonctionnalités métier exposées par l'API Go.
        </p>
      </div>
      <button class="button button-primary" @click="loadAll">Actualiser</button>
    </div>

    <LoadingState v-if="loading" message="Chargement du pilotage..." />
    <ErrorState v-else-if="error" :message="error" retry-label="Recharger" @retry="loadAll" />

    <template v-else>
      <div class="stat-grid">
        <StatCard
          v-for="metric in metrics"
          :key="metric.key"
          :label="metric.label"
          :value="metric.value"
          tone="green"
        />
      </div>

      <section class="surface-card section-card">
        <h3>Abonnements DM Plus</h3>
        <DataTable
          :columns="subscriptionColumns"
          :rows="subscriptions"
          row-key="id"
          empty-label="Aucun abonnement"
        />
      </section>

      <section class="surface-card section-card">
        <h3>Messagerie privée</h3>
        <DataTable
          :columns="messageColumns"
          :rows="messages"
          row-key="id"
          empty-label="Aucune conversation"
        />
      </section>

      <section class="commerce-grid">
        <div class="surface-card section-card">
          <h3>Offres négociées</h3>
          <DataTable :columns="offerColumns" :rows="commerce.offers" row-key="id" empty-label="Aucune offre" />
        </div>
        <div class="surface-card section-card">
          <h3>Ventes suivies</h3>
          <DataTable :columns="saleColumns" :rows="commerce.sales" row-key="id" empty-label="Aucune vente" />
        </div>
      </section>

      <section class="surface-card section-card">
        <div class="raw-header">
          <h3>Données métier</h3>
          <select v-model="rawTarget" @change="loadRaw">
            <option v-for="target in rawTargets" :key="target" :value="target">{{ target }}</option>
          </select>
        </div>
        <DataTable
          :columns="rawColumns"
          :rows="rawRows"
          row-key="id"
          empty-label="Aucune donnée"
        />
      </section>
    </template>
  </section>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from "vue";
import DataTable from "../components/DataTable.vue";
import ErrorState from "../components/ErrorState.vue";
import LoadingState from "../components/LoadingState.vue";
import StatCard from "../components/StatCard.vue";
import { adminApi } from "../services/api";

const loading = ref(false);
const error = ref("");
const metrics = ref([]);
const subscriptions = ref([]);
const messages = ref([]);
const commerce = reactive({ offers: [], sales: [] });
const rawTargets = ["users", "annonces", "events", "formations", "orders", "transactions", "reviews"];
const rawTarget = ref("users");
const rawRows = ref([]);

const subscriptionColumns = [
  { key: "id", label: "ID" },
  { key: "user", label: "Utilisateur" },
  { key: "email", label: "Email" },
  { key: "plan", label: "Plan" },
  { key: "price", label: "Prix" },
  { key: "status", label: "Statut" },
  { key: "date_fin", label: "Fin" }
];

const messageColumns = [
  { key: "id", label: "ID" },
  { key: "user_one", label: "Membre A" },
  { key: "user_two", label: "Membre B" },
  { key: "title", label: "Contexte" },
  { key: "updated_at", label: "MAJ" }
];

const offerColumns = [
  { key: "id", label: "ID" },
  { key: "conversation_id", label: "Conv." },
  { key: "buyer_id", label: "Acheteur" },
  { key: "seller_id", label: "Vendeur" },
  { key: "amount", label: "Montant" },
  { key: "status", label: "Statut" }
];

const saleColumns = [
  { key: "id", label: "ID" },
  { key: "offer_id", label: "Offre" },
  { key: "buyer_id", label: "Acheteur" },
  { key: "seller_id", label: "Vendeur" },
  { key: "amount", label: "Montant" },
  { key: "status", label: "Statut" },
  { key: "received_at", label: "Reception" },
  { key: "reviewed_at", label: "Avis" }
];

const rawColumns = computed(() => {
  const first = rawRows.value[0] || {};
  return Object.keys(first).map((key) => ({ key, label: key }));
});

async function loadRaw() {
  const response = await adminApi.rawDump(rawTarget.value);
  rawRows.value = response.items ?? [];
}

async function loadAll() {
  loading.value = true;
  error.value = "";
  try {
    const [overview, subs, convs, commerceResponse] = await Promise.all([
      adminApi.getSystemOverview(),
      adminApi.listSubscriptions(),
      adminApi.listMessages(),
      adminApi.listCommerce()
    ]);
    metrics.value = overview.metrics ?? [];
    subscriptions.value = subs.items ?? [];
    messages.value = convs.items ?? [];
    commerce.offers = commerceResponse.offers ?? [];
    commerce.sales = commerceResponse.sales ?? [];
    await loadRaw();
  } catch (err) {
    error.value = err.message || "Erreur de chargement du pilotage.";
  } finally {
    loading.value = false;
  }
}

onMounted(loadAll);
</script>

<style scoped>
.stack {
  display: grid;
  gap: 22px;
}

.section-card {
  padding: 24px;
}

.section-card h3 {
  margin: 0 0 16px;
}

.commerce-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 22px;
}

.raw-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 18px;
  margin-bottom: 16px;
}

.raw-header select {
  min-height: 42px;
  border: 1px solid rgba(141, 170, 152, 0.3);
  border-radius: 14px;
  padding: 0 14px;
  color: #eef7f1;
  background: rgba(255, 255, 255, 0.06);
}

@media (max-width: 1100px) {
  .commerce-grid {
    grid-template-columns: 1fr;
  }
}
</style>
