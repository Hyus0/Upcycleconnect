<template>
  <div class="layout-wrapper">    
    <div v-if="loading" class="state-card" style="margin: 40px auto; max-width: 600px; text-align: center;">
      <i class="ti ti-loader" style="font-size: 2rem; display: block; margin-bottom: 10px;"></i>
      Chargement de votre espace...
    </div>
    
    <main v-else>
      <DashboardPro v-if="isPremium" :user-id="currentUserId" />

      <DashboardClassique v-else :user-id="currentUserId" />
    </main>

  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue';
import SiteNavbar from "../../components/SiteNavbar.vue";
import SiteFooter from "../../components/SiteFooter.vue";

import DashboardClassique from './DashboardFreemium.vue';
import DashboardPro from './DashboardPremium.vue';

const API_URL = "/go";
const loading = ref(true);
const isPremium = ref(false);

const currentUserId = computed(() => Number(sessionStorage.getItem("userId")) || 0);
const isLoggedIn = computed(() => !!sessionStorage.getItem("userToken"));
const userName = computed(() => {
  const prenom = sessionStorage.getItem("userPrenom") || "";
  const nom = sessionStorage.getItem("userNom") || "";
  return (prenom || nom) ? `${prenom} ${nom}`.trim() : "Utilisateur";
});

const checkSubscription = async () => {
  if (!currentUserId.value) {
    loading.value = false;
    return;
  }
  
  try {
    const res = await fetch(`${API_URL}/users/${currentUserId.value}/abonnement`, {
      headers: { Authorization: `Bearer ${sessionStorage.getItem("userToken") || ""}` }
    });
    
    if (res.ok) {
      const data = await res.json();
      
      console.log("🕵️ Abonnement actif détecté par l'API :", data);
      
      const nomPlan = (data.plan_name || "").toLowerCase();
      const isPlanPro = data.plan_id === 2 || nomPlan.includes("pro") || nomPlan.includes("entreprise");
      
      isPremium.value = data.is_premium === true && isPlanPro;
    }
  } catch (error) {
    console.error("Erreur vérification abonnement:", error);
  } finally {
    loading.value = false;
  }
};

onMounted(() => {
  checkSubscription();
});
</script>

<style scoped>
.layout-wrapper {
  background: var(--bg-light, #f7f9f7);
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

main {
  flex: 1; 
}

.state-card {
  border: 1px dashed #cfe0d4;
  border-radius: 14px;
  padding: 40px 26px;
  color: var(--text-grey, #666);
  background: #fbfdfb;
}
</style>