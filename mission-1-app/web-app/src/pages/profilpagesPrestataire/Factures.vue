<template>
  <section class="factures-page">
    <header class="factures-hero">
      <p>ESPACE PRO</p>
      <h1>Factures</h1>
      <button type="button" @click="fetchFactures">Actualiser</button>
    </header>

    <div v-if="loading" class="factures-card">Chargement des factures...</div>
    <div v-else-if="error" class="factures-card error">{{ error }}</div>
    <div v-else-if="factures.length === 0" class="factures-card">
      <h2>Aucune facture pour le moment</h2>
      <p>Vos factures apparaitront ici apres un achat ou un abonnement.</p>
    </div>

    <div v-else class="factures-list">
      <article v-for="facture in factures" :key="facture.id" class="facture-row">
        <div>
          <span class="invoice-number">{{ facture.numero_facture }}</span>
          <h2>Commande #{{ facture.commande_id }}</h2>
          <p>{{ facture.date_transaction }} - {{ facture.statut_paiement }}</p>
        </div>
        <strong>{{ formatPrice(facture.montant_total) }}</strong>
        <div class="row-actions">
          <button type="button" @click="downloadFacture(facture.id)">Telecharger</button>
          <button type="button" class="ghost" @click="sendFacture(facture.id)">Envoyer par mail</button>
        </div>
      </article>
    </div>
  </section>
</template>

<script setup>
import { onMounted, ref } from "vue";

const API_URL = "/go";
const loading = ref(true);
const error = ref("");
const factures = ref([]);

const userId = () => Number(sessionStorage.getItem("id") || sessionStorage.getItem("userId")) || 0;

const formatPrice = (value) => new Intl.NumberFormat("fr-FR", {
  style: "currency",
  currency: "EUR"
}).format(Number(value) || 0);

const fetchFactures = async () => {
  loading.value = true;
  error.value = "";
  try {
    const res = await fetch(`${API_URL}/users/${userId()}/factures`, {
      headers: { Authorization: `Bearer ${sessionStorage.getItem("userToken")}` }
    });
    if (!res.ok) throw new Error(await res.text());
    factures.value = await res.json() || [];
  } catch (err) {
    error.value = err.message || "Impossible de charger les factures.";
  } finally {
    loading.value = false;
  }
};

const downloadFacture = (factureId) => {
  window.open(`${API_URL}/users/${userId()}/factures/${factureId}/download`, "_blank");
};

const sendFacture = async (factureId) => {
  const res = await fetch(`${API_URL}/users/${userId()}/factures/${factureId}/send`, {
    method: "POST",
    headers: { Authorization: `Bearer ${sessionStorage.getItem("userToken")}` }
  });
  const payload = res.ok ? await res.json() : { message: await res.text() };
  alert(payload.message || "Demande traitee.");
};

onMounted(fetchFactures);
</script>

<style scoped>
.factures-page {
  padding: 32px;
  font-family: "Syne", sans-serif;
}
.factures-hero {
  display: flex;
  align-items: end;
  justify-content: space-between;
  gap: 24px;
  margin-bottom: 28px;
}
.factures-hero p {
  color: #2f8f58;
  letter-spacing: 4px;
  font-size: 0.75rem;
  font-weight: 900;
  margin: 0 0 8px;
}
.factures-hero h1 {
  margin: 0;
  color: #16221c;
  font-size: clamp(2rem, 6vw, 4.5rem);
}
button {
  border: 0;
  border-radius: 14px;
  padding: 13px 18px;
  background: #2f8f58;
  color: #fff;
  font-weight: 800;
  cursor: pointer;
}
button.ghost {
  background: #e9f3ee;
  color: #1f6b43;
}
.factures-card,
.facture-row {
  background: #fff;
  border: 1px solid #dce9e1;
  border-radius: 22px;
  padding: 26px;
  box-shadow: 0 16px 42px rgba(24, 56, 39, 0.06);
}
.factures-card.error {
  border-color: #f3b7ac;
  color: #b3261e;
}
.factures-list {
  display: grid;
  gap: 16px;
}
.facture-row {
  display: grid;
  grid-template-columns: 1fr auto auto;
  align-items: center;
  gap: 20px;
}
.invoice-number {
  color: #2f8f58;
  font-size: 0.82rem;
  font-weight: 900;
  letter-spacing: 1px;
}
.facture-row h2 {
  margin: 6px 0;
  color: #17201b;
}
.facture-row p {
  margin: 0;
  color: #718176;
}
.row-actions {
  display: flex;
  gap: 10px;
}
@media (max-width: 760px) {
  .factures-hero,
  .facture-row {
    display: flex;
    flex-direction: column;
    align-items: flex-start;
  }
  .row-actions {
    flex-wrap: wrap;
  }
}
</style>
