<template>
    <main class="public-dashboard">
        <SiteNavbar
            :is-authenticated="isLoggedIn"
            :user-name="userName"
            variant="public"
        />

        <header class="content-header">
            <div class="header-left">
                <p class="sidebar__category2">
                    ACCUEIL > MON COMPTE > ESPACE PRO
                </p>
                <h1 class="hero-title1">Mon Espace Pro</h1>
                <p class="classic-text">
                    Gérez votre accès aux outils avancés, vos collectes et vos
                    campagnes.
                </p>
            </div>
            <button
                class="btn-secondary"
                type="button"
                @click="$router.push('/profil')"
            >
                🠔 Retour au profil
            </button>
        </header>

        <div
            v-if="loading"
            class="state-card"
            style="margin-top: 24px; text-align: center"
        >
            Chargement de votre espace...
        </div>        

        <section v-else class="section-container" style="margin-top: 24px">
                
            <div
                v-if="!subscription.is_premium"
                class="pricing-container text-center"
            >
                <div class="section-header justify-center" style="margin-bottom: 2rem">
                    <div>
                        <h2 style="font-size: 1.8rem; margin-bottom: 8px">
                            Passez à la vitesse supérieure
                        </h2>
                        <p class="classic-text">
                            Débloquez des outils puissants pour votre activité.
                        </p>
                    </div>
                </div>

                <div class="pricing-grid">
                    <div
                        v-for="plan in sortedPlans"
                        :key="plan.id"
                        class="card pricing-card"
                        :class="plan.prix_ht > 0 ? 'premium-card' : 'card--white'"
                        :style="plan.color && plan.color.toLowerCase() !== 'ffffff' ? { backgroundColor: '#' + plan.color, borderColor: '#' + plan.color } : {}"
                    >
                        <div v-if="plan.prix_ht > 0" class="tag-score popular-tag" :style="subscription.is_premium ? { background: 'white', color: '#' + plan.color } : { background: '#ffb300', color: '#333' }">
                            {{ subscription.is_premium ? 'ACTIF' : 'RECOMMANDÉ' }}
                        </div>
                        
                        <h3 class="pricing-title" :class="{ 'text-white': plan.color.toLowerCase() !== 'ffffff' }">
                            {{ plan.nom }}
                        </h3>
                        
                        <p class="pricing-price" :class="{ 'text-white': plan.color.toLowerCase() !== 'ffffff' }">
                            {{ Number(plan.prix_ht).toFixed(2) }} €
                            <span :class="{ 'text-white-50': plan.color.toLowerCase() !== 'ffffff' }">/mois</span>
                        </p>
                        
                        <ul class="pricing-features" :class="{ 'text-white': plan.color.toLowerCase() !== 'ffffff' }">
                            <li v-for="avantage in plan.avantages" :key="avantage.id" :class="{ 'locked': !avantage.disponible }">
                                <template v-if="avantage.disponible">
                                    <Check :size="18" :class="plan.color.toLowerCase() !== 'ffffff' ? 'text-white' : 'text-green-500'" />
                                    <strong v-if="avantage.nom.toLowerCase().includes('vitrine')">{{ avantage.nom }}</strong>
                                    <span v-else>{{ avantage.nom }}</span>
                                </template>
                                <template v-else>
                                    🔒 {{ avantage.nom }}
                                </template>
                            </li>
                        </ul>

                        <button
                            v-if="plan.prix_ht == 0 && subscription.is_premium"
                            class="btn-danger w-full"
                            style="width: 100%; margin-top: auto; justify-content: center;"
                            :disabled="processing"
                            @click="resilierAbonnement"
                        >
                            {{ processing ? "Traitement..." : "Résilier et repasser au gratuit" }}
                        </button>

                        <button
                            v-else-if="plan.prix_ht == 0"
                            class="btn-secondary w-full"
                            disabled
                            style="width: 100%; margin-top: auto;"
                        >
                            Plan actuel
                        </button>

                        <button
                            v-else-if="plan.prix_ht > 0 && subscription.is_premium && subscription.plan_id === plan.id"
                            class="btn-secondary w-full"
                            disabled
                            style="width: 100%; margin-top: auto; color: #2d7a4f; background: white; opacity: 0.9;"
                        >
                            Plan actuel
                        </button>

                        <div v-else style="width: 100%; margin-top: auto; display: flex; flex-direction: column; gap: 8px;">
                            <button
                                class="btn-main-action w-full justify-center btn-premium"
                                style="width: 100%; justify-content: center;"
                                :disabled="processing"
                                @click="goToPaiement(plan)"
                            >
                                S'abonner maintenant
                            </button>
                            
                            <button
                                class="btn-test w-full justify-center"
                                style="width: 100%; justify-content: center; padding: 8px; font-size: 0.85rem;"
                                :disabled="processing"
                                @click="souscrireDirectement(plan)"
                            >
                                {{ processing ? "Activation..." : "Test: Activer sans payer" }}
                            </button>
                        </div>
                    </div>
                </div>
            </div>

            <div v-else>

                <div class="pricing-container text-center" style="margin-bottom: 4rem">
                    <div class="section-header justify-center" style="margin-bottom: 2rem">
                        <div>
                            <h2 style="font-size: 1.8rem; margin-bottom: 8px">
                                Gestion de votre abonnement
                            </h2>
                            <p class="classic-text" style="font-size: 1.1rem">
                                Votre abonnement est actuellement
                                <strong>{{ subscription.statut.toLowerCase() }}</strong>
                                Il est valable jusqu'au
                                <strong>{{ formatDate(subscription.date_fin) || "Facturation continue" }}</strong>.
                            </p>
                        </div>
                    </div>

                    <div class="pricing-grid">
                        <div
                            v-for="plan in sortedPlans"
                            :key="plan.id"
                            class="card pricing-card"
                            :class="plan.prix_ht > 0 ? 'premium-card' : 'card--white'"
                            :style="plan.color && plan.color.toLowerCase() !== 'ffffff' ? { backgroundColor: '#' + plan.color, borderColor: '#' + plan.color } : {}"
                        >
                            <div v-if="plan.prix_ht > 0" class="tag-score popular-tag" :style="subscription.is_premium ? { background: 'white', color: '#' + plan.color } : { background: '#ffb300', color: '#333' }">
                                {{ subscription.is_premium ? 'ACTIF' : 'RECOMMANDÉ' }}
                            </div>
                            
                            <h3 class="pricing-title" :class="{ 'text-white': plan.color.toLowerCase() !== 'ffffff' }">
                                {{ plan.nom }}
                            </h3>
                            
                            <p class="pricing-price" :class="{ 'text-white': plan.color.toLowerCase() !== 'ffffff' }">
                                {{ Number(plan.prix_ht).toFixed(2) }} €
                                <span :class="{ 'text-white-50': plan.color.toLowerCase() !== 'ffffff' }">/mois</span>
                            </p>
                            
                            <ul class="pricing-features" :class="{ 'text-white': plan.color.toLowerCase() !== 'ffffff' }">
                                <li v-for="avantage in plan.avantages" :key="avantage.id" :class="{ 'locked': !avantage.disponible }">
                                    <template v-if="avantage.disponible">
                                        <Check :size="18" :class="plan.color.toLowerCase() !== 'ffffff' ? 'text-white' : 'text-green-500'" />
                                        <strong v-if="avantage.nom.toLowerCase().includes('vitrine')">{{ avantage.nom }}</strong>
                                        <span v-else>{{ avantage.nom }}</span>
                                    </template>
                                    <template v-else>
                                        🔒 {{ avantage.nom }}
                                    </template>
                                </li>
                            </ul>

                            <button
                                v-if="plan.prix_ht == 0 && subscription.is_premium"
                                class="btn-danger w-full"
                                style="width: 100%; margin-top: auto; justify-content: center;"
                                :disabled="processing"
                                @click="resilierAbonnement"
                            >
                                {{ processing ? "Traitement..." : "Résilier et repasser au gratuit" }}
                            </button>

                            <button
                                v-else-if="plan.prix_ht == 0"
                                class="btn-secondary w-full"
                                disabled
                                style="width: 100%; margin-top: auto"
                            >
                                Plan actuel
                            </button>

                            <button
                                v-else-if="plan.prix_ht > 0 && subscription.is_premium && subscription.plan_id === plan.id"
                                class="btn-secondary w-full"
                                disabled
                                style="width: 100%; margin-top: auto; color: #2d7a4f; background: white; opacity: 0.9;"
                            >
                                Plan actuel
                            </button>

                            <div v-else style="width: 100%; margin-top: auto; display: flex; flex-direction: column; gap: 8px;">
                                <button
                                    class="btn-main-action w-full justify-center btn-premium"
                                    style="width: 100%; justify-content: center;"
                                    :disabled="processing"
                                    @click="goToPaiement(plan)"
                                >
                                    S'abonner maintenant
                                </button>
                                
                                <button
                                    class="btn-test w-full justify-center"
                                    style="width: 100%; justify-content: center; padding: 8px; font-size: 0.85rem;"
                                    :disabled="processing"
                                    @click="souscrireDirectement(plan)"
                                >
                                    {{ processing ? "Activation..." : "Test: Activer sans payer" }}
                                </button>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </section>
    </main>
    <SiteFooter />
</template>

<script setup>
import { ref, onMounted, computed } from "vue";
import { useRouter } from "vue-router";
import SiteNavbar from "../components/SiteNavbar.vue";
import SiteFooter from "../components/SiteFooter.vue";
import { Check } from "lucide-vue-next";

const router = useRouter();
const API_URL = "/go";

const loading = ref(true);
const processing = ref(false);

const activeTab = ref("dashboard");
const stats = ref({});
const priorityAlerts = ref([]);
const campaigns = ref([]);

const subscription = ref({ is_premium: false, statut: "" });
const AllAbonnements = ref([]);

const currentUserId = computed(
    () => Number(sessionStorage.getItem("userId")) || 0,
);
const isLoggedIn = computed(() => !!sessionStorage.getItem("userToken"));
const userName = computed(() => {
    const prenom = sessionStorage.getItem("userPrenom") || "";
    const nom = sessionStorage.getItem("userNom") || "";
    return prenom || nom ? `${prenom} ${nom}`.trim() : "Utilisateur";
});

const sortedPlans = computed(() => {
    return AllAbonnements.value
        .filter(plan => plan.id !== 1) 
        .sort((a, b) => a.prix_ht - b.prix_ht);
});

const getHeaders = () => ({
    "Content-Type": "application/json",
    Authorization: sessionStorage.getItem("userToken") || "",
});

const formatDate = (dateStr) => {
    if (!dateStr) return "";
    const date = new Date(dateStr.replace(/Z$/, ""));
    if (Number.isNaN(date.getTime())) return "";
    return new Intl.DateTimeFormat("fr-FR", {
        day: "2-digit",
        month: "long",
        year: "numeric",
    }).format(date);
};

const fetchAllAbonnements = async () => {
    try {
        const res = await fetch(`${API_URL}/abonnements`);
        if (res.ok) {
            AllAbonnements.value = await res.json();
        } else {
            console.error("Erreur lors de la récupération des plans");
        }
    } catch (err) {
        console.error("Erreur réseau:", err);
    }
};

const souscrireDirectement = async (plan) => {
    if (!currentUserId.value) {
        alert("Veuillez vous connecter pour vous abonner.");
        router.push("/connexion");
        return;
    }

    processing.value = true;
    try {
        const res = await fetch(
            `${API_URL}/users/${currentUserId.value}/abonnement/souscrire`,
            {
                method: "POST",
                headers: getHeaders(),
                body: JSON.stringify({
                    plan_id: plan.id, 
                    stripe_payment_id: "bypass_stripe_manuel",
                }),
            },
        );

        if (res.ok) {
            alert(`Félicitations, votre abonnement ${plan.nom} est activé !`);
            await fetchAbonnement();
        } else {
            const errMsg = await res.text();
            alert("Erreur lors de l'activation : " + errMsg);
        }
    } catch (error) {
        console.error("Erreur abonnement:", error);
        alert("Impossible de contacter le serveur.");
    } finally {
        processing.value = false;
    }
};

const loadProData = async () => {
    stats.value = {
        materiaux_sauves_kg: 1250,
        co2_evite_kg: 450,
        top_materiau: "Palettes Bois",
        projets_actifs: 4,
    };

    priorityAlerts.value = [
        {
            id: 1,
            materiau: "Bois Massif (50kg)",
            localisation: "Déchèterie Nord",
            temps: "Il y a 2h",
        },
        {
            id: 2,
            materiau: "Chutes de cuir",
            localisation: "Atelier Textile Centre",
            temps: "Il y a 5h",
        },
    ];

    campaigns.value = [
        {
            id: 1,
            titre: "Table basse industrielle",
            budget: 150,
            vues: 4230,
            statut: "Active",
        },
    ];
};

const fetchAbonnement = async () => {
    if (!currentUserId.value) {
        loading.value = false;
        return;
    }
    loading.value = true;
    try {
        const res = await fetch(
            `${API_URL}/users/${currentUserId.value}/abonnement`,
            {
                headers: getHeaders(),
            },
        );
        if (res.ok) {
            subscription.value = await res.json();
            if (subscription.value.is_premium) {
                await loadProData();
            }
        }
    } catch (error) {
        console.error("Erreur de chargement de l'abonnement:", error);
    } finally {
        loading.value = false;
    }
};

const goToPaiement = (plan) => {
    if (!currentUserId.value) {
        alert("Veuillez vous connecter pour vous abonner.");
        router.push("/connexion");
        return;
    }

    router.push({
        path: "/paiement",
        query: {
            plan_id: plan.id,
            prix: plan.prix_ht,
            nom_plan: plan.nom,
        },
    });
};

const resilierAbonnement = async () => {
    if (
        !confirm(
            "Êtes-vous sûr de vouloir résilier ? Vous garderez l'accès jusqu'à la fin du mois en cours.",
        )
    )
        return;

    processing.value = true;
    try {
        const res = await fetch(
            `${API_URL}/users/${currentUserId.value}/abonnement/resilier`,
            {
                method: "POST",
                headers: getHeaders(),
            },
        );

        if (res.ok) {
            alert("Votre abonnement a été résilié.");
            await fetchAbonnement();
        } else {
            alert("Erreur lors de la résiliation.");
        }
    } catch (error) {
        console.error("Erreur de résiliation:", error);
    } finally {
        processing.value = false;
    }
};

onMounted(async () => {
    await fetchAllAbonnements();
    await fetchAbonnement();
});
</script>

<style scoped>
.public-dashboard {
    min-height: 100vh;
    padding: 20px;
    background: var(--bg-light, #f7f9f7);
}

.content-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-end;
    margin-bottom: 2rem;
    padding-bottom: 1.5rem;
    border-bottom: 1px solid #e5ede7;
}

.sidebar__category2 {
    margin: 0;
    color: #a0ada7;
    font-family: "Space Mono", monospace;
    font-size: 0.65rem;
    letter-spacing: 1px;
    text-transform: uppercase;
}

.hero-title1 {
    font-size: 2.2rem;
    font-weight: 900;
    color: #122018;
    margin: 0;
}

.classic-text {
    font-size: 0.95rem;
    color: #63746a;
    margin-top: 8px;
    line-height: 1.5;
}

.btn-secondary {
    padding: 8px 16px;
    border-radius: 10px;
    border: 1px solid #ddd;
    background: white;
    cursor: pointer;
    font-weight: bold;
}

.btn-secondary:disabled {
    background: #f5f5f5;
    color: #999;
    border-color: #eee;
    cursor: not-allowed;
}

.btn-test {
    background-color: #f3f4f6;
    color: #374151;
    border: 1px solid #d1d5db;
    border-radius: 10px;
    font-weight: 700;
    cursor: pointer;
    transition: background 0.2s;
}

.btn-test:hover:not(:disabled) {
    background-color: #e5e7eb;
}

.section-container {
    background: white;
    padding: 2rem;
    border-radius: 16px;
    border: 1px solid #eee;
    margin-bottom: 2rem;
}

.state-card {
    border: 1px dashed #cfe0d4;
    border-radius: 14px;
    padding: 26px;
    color: var(--text-grey, #666);
    background: #fbfdfb;
}

.pricing-container {
    display: flex;
    flex-direction: column;
    align-items: center;
}

.pricing-grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 2rem;
    max-width: 900px;
    width: 100%;
}

@media (max-width: 768px) {
    .pricing-grid {
        grid-template-columns: 1fr;
    }
}

.pricing-card {
    padding: 2.5rem;
    border-radius: 16px;
    border: 1px solid #eee;
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.02);
    display: flex;
    flex-direction: column;
    position: relative;
    text-align: left;
}

.premium-card {
    background: #2d7a4f;
    color: white;
    border: none;
    transform: scale(1.02);
    box-shadow: 0 24px 60px rgba(16, 32, 24, 0.25);
}

.pricing-title {
    font-size: 1.3rem;
    font-weight: 800;
    margin: 0;
    font-family: "Syne", sans-serif;
}

.pricing-price {
    font-size: 2.8rem;
    font-weight: 900;
    margin: 1rem 0;
}

.pricing-price span {
    font-size: 1.1rem;
    font-weight: normal;
}

.pricing-features {
    list-style: none;
    padding: 0;
    margin: 2rem 0;
}

.pricing-features li {
    display: flex;
    align-items: center;
    gap: 12px;
    margin-bottom: 1rem;
    font-size: 0.95rem;
    font-weight: 600;
}

.pricing-features .locked {
    color: #999;
    text-decoration: line-through;
}

.popular-tag {
    position: absolute;
    top: 0;
    right: 0;
    border-radius: 0 16px 0 12px;
    background: #ffb300;
    color: #333;
    font-weight: bold;
    font-size: 0.75rem;
    padding: 6px 12px;
}

.text-white {
    color: white !important;
}
.text-white-50 {
    color: rgba(255, 255, 255, 0.7) !important;
}
.text-green-500 {
    color: #2d7a4f;
}

.tabs-container {
    display: flex;
    gap: 1rem;
    margin-bottom: 2rem;
    border-bottom: 2px solid #eee;
    padding-bottom: 10px;
}

.tab-btn {
    background: transparent;
    border: none;
    font-size: 1.1rem;
    font-weight: bold;
    color: #666;
    padding: 8px 16px;
    cursor: pointer;
    border-radius: 8px;
    transition: all 0.2s;
}

.tab-btn:hover {
    background: #eaf4ed;
    color: #2d7a4f;
}
.tab-btn.active {
    background: #2d7a4f;
    color: white;
}

.stats-grid {
    display: grid;
    gap: 1.5rem;
    grid-template-columns: 1.4fr 0.8fr 0.8fr;
}

@media (max-width: 920px) {
    .stats-grid {
        grid-template-columns: 1fr;
    }
}

.card {
    background: white;
    padding: 1.5rem;
    border-radius: 12px;
    border: 1px solid #eee;
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.02);
}

.card--score {
    background: #2d7a4f;
    color: white;
    border: none;
}
.tag-score {
    background: rgba(255, 255, 255, 0.2);
    color: white;
    padding: 4px 10px;
    border-radius: 6px;
    font-size: 0.75rem;
    font-weight: 800;
    display: inline-block;
    margin-bottom: 1rem;
}

.score-value {
    font-size: 2.8rem;
    font-weight: 900;
    line-height: 1;
    margin-bottom: 0.5rem;
}

.score-value span {
    font-size: 1.2rem;
    font-weight: normal;
    opacity: 0.8;
}
.card-num,
.card-num2 {
    font-size: 2.5rem;
    font-weight: bold;
    line-height: 1;
    margin-bottom: 0.5rem;
}

.score-footer {
    display: flex;
    justify-content: space-between;
    border-top: 1px solid rgba(255, 255, 255, 0.2);
    padding-top: 1rem;
    margin-top: 1rem;
}

.badge {
    padding: 4px 10px;
    border-radius: 20px;
    font-size: 0.8rem;
    font-weight: bold;
    display: inline-block;
}

.badge--green {
    background: #eaf4ed;
    color: #2d7a4f;
}
.badge--orange {
    background: #fff3cd;
    color: #856404;
}

.data-table {
    width: 100%;
    border-collapse: collapse;
}

.data-table th,
.data-table td {
    padding: 12px;
    text-align: left;
    border-bottom: 1px solid #eee;
}

.status-logistique {
    background: #f0f0f0;
    color: #555;
    padding: 4px 10px;
    border-radius: 6px;
    font-size: 0.75rem;
    font-weight: bold;
    text-transform: uppercase;
}

.btn-main-action {
    display: inline-flex;
    align-items: center;
    background: #2d7a4f;
    color: white;
    padding: 10px 20px;
    border-radius: 10px;
    text-decoration: none;
    font-weight: bold;
    border: none;
    cursor: pointer;
}
.btn-main-action:hover:not(:disabled) {
    background: #1b4d31;
}

.btn-premium {
    background: white;
    color: #2d7a4f;
}
.btn-premium:hover:not(:disabled) {
    background: #f0f0f0;
}

.btn-danger {
    background: transparent;
    color: #c94b4b;
    border: 1px solid #f7d7d7;
    border-radius: 12px;
    padding: 10px 18px;
    font-size: 0.95rem;
    font-family: "Syne", sans-serif;
    font-weight: 700;
    cursor: pointer;
    transition: 0.2s;
}
.btn-danger:hover:not(:disabled) {
    background: #fceaea;
}
.btn-danger:disabled,
.btn-main-action:disabled {
    opacity: 0.6;
    cursor: not-allowed;
}
</style>