<template>
    <main class="page-main-content">
        <header class="content-header">
            <div class="header-left">
                <p class="sidebar__category2">
                    ESPACE PRESTATAIRE > LOGISTIQUE
                </p>
                <h1 class="hero-title1">MES RÉCUPÉRATIONS {{ id }}</h1>
                <p class="classic-text">
                    Suivez l'état de vos acquisitions et organisez vos tournées
                    de collecte.
                </p>
            </div>
            <button class="btn-main-action" @click="$router.push('/catalogue')">
                Voir le catalogue
            </button>
        </header>

        <div v-if="loading" class="loading-state">
            Chargement de vos acquisitions...
        </div>

        <div v-else class="section-container">
            <div class="dash-block">
                <h2 class="block-title">
                    EN ATTENTE DU VENDEUR
                    <span class="badge">{{ enAttente.length }}</span>
                </h2>
                <table class="data-table">
                    <thead>
                        <tr>
                            <th>OBJET</th>
                            <th>TYPE</th>
                            <th>DATE D'ACHAT</th>
                            <th>STATUT ACTUEL</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="item in enAttente" :key="item.id">
                            <td>
                                <strong>{{ item.titre }}</strong>
                            </td>
                            <td>
                                <span
                                    :class="
                                        item.type === 'Don'
                                            ? 'tag-don'
                                            : 'tag-vente'
                                    "
                                >
                                    {{
                                        item.type === "Don"
                                            ? "🎁 DON"
                                            : "💰 ACHAT"
                                    }}
                                </span>
                            </td>
                            <td>{{ formatDate(item.date_achat) }}</td>
                            <td>
                                <span class="status-pill reserve">
                                    {{
                                        item.statut === "Paye"
                                            ? "Le vendeur doit réserver un casier"
                                            : "Casier réservé, en attente de dépôt"
                                    }}
                                </span>
                            </td>
                        </tr>
                    </tbody>
                </table>
                <p v-if="enAttente.length === 0" class="empty-text">
                    Vous n'avez aucun objet en attente de la part des vendeurs.
                </p>
            </div>

            <div class="dash-block ">
                <h2 class="block-title title-urgent">
                    DISPONIBLES EN BORNE
                    <span class="badge badge-urgent">{{
                        pretsAuRetrait.length
                    }}</span>
                </h2>
                <table class="data-table">
                    <thead>
                        <tr>
                            <th>OBJET</th>
                            <th>SITE DE COLLECTE</th>
                            <th>CODE DE RETRAIT (PIN)</th>
                            <th class="text-right">ACTION</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="item in pretsAuRetrait" :key="item.id">
                            <td>
                                <strong>{{ item.titre }}</strong
                                ><br />
                                <small
                                    >Poids :
                                    {{ item.poids_estime_kg }} kg</small
                                >
                            </td>
                            <td>
                                <div
                                    class="site-cell clickable"
                                    @click="fetchAndOpenSite(item.id_site)"
                                >
                                    <strong>{{
                                        sitesCache[item.id_site] ||
                                        "Chargement..."
                                    }}</strong
                                    ><br />
                                    <small
                                        >Casier n°{{
                                            item.id_casier || "..."
                                        }}
                                        (Détails ↗)</small
                                    >
                                </div>
                            </td>
                            <td>
                                <div
                                    class="token-cell clickable"
                                    @click.stop.prevent="
                                        showToken(item.code_barre_retrait)
                                    "
                                >
                                    <span class="pin-badge">
                                        {{
                                            truncateToken(
                                                item.code_barre_retrait,
                                            )
                                        }}
                                    </span>
                                </div>
                            </td>
                            <td class="text-right">
                                <button
                                    class="btn-plan"
                                    @click="
                                        $router.push('/profil/recuperer-objet')
                                    "
                                >
                                    Simuler le retrait 🏃‍♂️
                                </button>
                            </td>
                        </tr>
                    </tbody>
                </table>
                <p v-if="pretsAuRetrait.length === 0" class="empty-text">
                    Aucun objet en attente de récupération physique.
                </p>
            </div>

            <div class="dash-block history-bg">
                <div class="flex-between">
                    <h2 class="block-title">
                        HISTORIQUE (OBJETS RÉCUPÉRÉS)
                    </h2>
                    <div class="search-mini">
                        <input
                            type="text"
                            v-model="searchQuery"
                            placeholder="Rechercher..."
                        />
                    </div>
                </div>
                <table class="data-table">
                    <thead>
                        <tr>
                            <th>OBJET</th>
                            <th>RÉCUPÉRÉ AU SITE</th>
                            <th class="text-center">DATE DE RETRAIT</th>
                            <th class="text-center">FACTURE</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="item in historiqueFiltre" :key="item.id">
                            <td>
                                <strong>{{ item.titre }}</strong>
                            </td>
                            <td>
                                <small>{{
                                    sitesCache[item.id_site] || "Site inconnu"
                                }}</small>
                            </td>
                            <td class="text-center small-date">
                                {{
                                    formatDate(item.date_recuperation_effective)
                                }}
                            </td>
                            <td class="text-center">
                                <button
                                    class="btn-facture"
                                    @click="telechargerFacture(item.id)"
                                    title="Télécharger la facture PDF"
                                >
                                    Télécharger PDF
                                </button>
                            </td>
                        </tr>
                    </tbody>
                </table>
                <p v-if="historiqueFiltre.length === 0" class="empty-text">
                    Aucun historique correspondant.
                </p>
            </div>
        </div>

        <div
            v-if="selectedSite"
            class="modal-overlay"
            @click.self="selectedSite = null"
        >
            <div class="site-modal">
                <div class="modal-header">
                    <h3>Détails du Point de Collecte</h3>
                    <button class="close-btn" @click="selectedSite = null">
                        ×
                    </button>
                </div>
                <div class="modal-content">
                    <div class="info-row">
                        <label>NOM</label>
                        <p>{{ selectedSite.nom }}</p>
                    </div>
                    <div class="info-row">
                        <label>ADRESSE</label>
                        <p>
                            {{ selectedSite.adresse }}<br />
                            {{ selectedSite.code_postal }}
                            {{ selectedSite.ville }}
                        </p>
                    </div>
                    <div class="info-row">
                        <label>TÉLÉPHONE</label>
                        <p>{{ selectedSite.telephone || "Non renseigné" }}</p>
                    </div>
                </div>
                <button
                    class="btn-main-action modal-close-btn"
                    @click="selectedSite = null"
                >
                    Fermer
                </button>
            </div>
        </div>

        <div
            v-if="selectedToken"
            class="modal-overlay"
            @click.self="selectedToken = null"
        >
            <div class="site-modal">
                <div class="modal-header">
                    <h3>Votre Code de Retrait</h3>
                    <button class="close-btn" @click="selectedToken = null">
                        ×
                    </button>
                </div>
                <div class="modal-content">
                    <p class="info-text">
                        Présentez ce QR Code ou le code unique au niveau du conteneur sécurisé pour déverrouiller le
                        casier : :
                    </p>
                
                    <div class="qr-wrapper">
                        <QrcodeVue
                            :value="selectedToken"
                            :size="180"
                            level="H"
                            background="#ffffff"
                            foreground="#1d3528"
                        />
                    </div>
                </div>
                <div class="token-box">
                    {{ selectedToken }}
                </div>
                <button
                    class="btn-main-action modal-close-btn"
                    @click="selectedToken = null"
                >
                    Fermer
                </button>
            </div>
        </div>
    </main>
</template>

<script setup>
import { ref, onMounted, computed } from "vue";
import { useRouter } from "vue-router";
import QrcodeVue from "qrcode.vue";

const router = useRouter();

const isLoggedIn = computed(() => !!sessionStorage.getItem("userToken"));
const userName = computed(
    () => sessionStorage.getItem("userPrenom") || "Prestataire",
);
const userRole = computed(
    () => sessionStorage.getItem("userRole") || "Prestataire",
);

const achats = ref([]);
const sitesCache = ref({});
const loading = ref(true);
const searchQuery = ref("");
const selectedSite = ref(null);
const selectedToken = ref(null);

const enAttente = computed(() =>
    achats.value.filter((a) => a.statut === "Paye" || a.statut === "Reserve"),
);

const pretsAuRetrait = computed(() =>
    achats.value.filter((a) => a.statut === "Depose"),
);

const historiqueFiltre = computed(() => {
    const list = achats.value.filter((a) => a.statut === "Recupere");
    return list.filter((a) =>
        a.titre.toLowerCase().includes(searchQuery.value.toLowerCase()),
    );
});

const telechargerFacture = (idAnnonce) => {
    alert(
        `Génération de la facture PDF pour l'achat #${idAnnonce} en cours...\n(À connecter avec le backend Go 🚀)`,
    );
};

const formatDate = (d) => {
    if (!d || d.startsWith("0001") || d === "null") return "---";
    return new Date(d).toLocaleDateString("fr-FR", {
        day: "numeric",
        month: "short",
        year: "numeric",
    });
};

const fetchSitesCache = async () => {
    try {
        const res = await fetch("/go/sites");
        if (res.ok) {
            const data = await res.json();
            data.forEach((s) => {
                sitesCache.value[s.id] = s.nom;
            });
        }
    } catch (e) {
        console.error(e);
    }
};

const truncateToken = (token) => {
    if (!token) return "";
    return token.length > 10 ? token.substring(0, 10) + "..." : token;
};

const showToken = (token) => {
    selectedToken.value = token;
};

const fetchAchats = async () => {
    const id = sessionStorage.getItem("userId");
    const token = sessionStorage.getItem("userToken");
    loading.value = true;
    try {
        const res = await fetch(`/go/users/${id}/achats`, {
            headers: { Authorization: token },
        });
        if (res.ok) achats.value = await res.json();
    } catch (e) {
        console.error("Erreur de récupération des achats", e);
    } finally {
        loading.value = false;
    }
};

const fetchAndOpenSite = async (siteId) => {
    if (!siteId) return;
    const token = sessionStorage.getItem("userToken");
    try {
        const res = await fetch(`/go/site/${siteId}`, {
            headers: { Authorization: token },
        });
        if (res.ok) {
            selectedSite.value = await res.json();
        }
    } catch (e) {
        console.error(e);
    }
};

onMounted(() => {
    fetchSitesCache();
    fetchAchats();
});
</script>

<style scoped>
.page-main-content {
    min-height: 100vh;
    padding: 20px;
    background: #f7f9f7;
    max-width: 1600px;
    margin: 0 auto;
}

.content-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-end;
    margin-bottom: 2rem;
    padding-bottom: 1.5rem;
    border-bottom: 1px solid #eee;
}

.header-left {
    display: flex;
    flex-direction: column;
}

.qr-wrapper {
    display: flex;
    justify-content: center;
    margin: 20px 0;
    padding: 15px;
    background: #fff;
    border-radius: 16px;
    border: 1px solid #e8ece9;
}

.qr-wrapper canvas,
.qr-wrapper svg {
    border-radius: 10px;
}

.sidebar__category2 {
    font-size: 0.65rem;
    color: #8fa396;
    letter-spacing: 1px;
    margin: 0 0 0.5rem 0;
    text-transform: uppercase;
}

.hero-title1 {
    font-size: 2rem;
    font-weight: 800;
    margin: 1.5rem 0 0.5rem;
    color: #1a1a1a;
}

.classic-text {
    color: #666;
    margin: 0;
}

.btn-main-action {
    background: #2d7a4f;
    color: white;
    padding: 10px 20px;
    border-radius: 10px;
    font-weight: bold;
    border: none;
    cursor: pointer;
}

.dash-block {
    background: white;
    border: 1px solid #eee;
    border-radius: 16px;
    padding: 1.5rem;
    margin-bottom: 2rem;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.02);
}



.block-title {
    font-size: 0.8rem;
    font-weight: 800;
    color: #999;
    margin-bottom: 1.5rem;
    display: flex;
    align-items: center;
    gap: 10px;
}

.title-urgent {
    color: #2d7a4f;
}

.badge {
    background: #f0f0f0;
    color: #666;
    padding: 2px 8px;
    border-radius: 10px;
    font-size: 0.7rem;
}

.badge-urgent {
    background: #e8f5e9;
    color: #2e7d32;
}

.data-table {
    width: 100%;
    border-collapse: collapse;
}
.data-table th {
    text-align: left;
    font-size: 0.7rem;
    color: #bbb;
    padding-bottom: 12px;
    border-bottom: 1px solid #eee;
}
.data-table td {
    padding: 14px 0;
    border-bottom: 1px solid #f9f9f9;
    font-size: 0.9rem;
    vertical-align: middle;
}

.text-right {
    text-align: right !important;
}
.text-center {
    text-align: center !important;
}

.pin-badge {
    background: #e8f5e9;
    color: #2d6a4f;
    font-weight: 900;
    padding: 6px 10px;
    border-radius: 6px;
    letter-spacing: 1px;
    font-family: monospace;
    font-size: 1rem;
}

.status-pill {
    font-size: 0.75rem;
    font-weight: bold;
    padding: 6px 10px;
    border-radius: 6px;
}
.status-pill.reserve {
    background: #fff3e0;
    color: #e65100;
}

.btn-plan {
    background: #2d6a4f;
    color: white;
    border: none;
    padding: 8px 16px;
    border-radius: 8px;
    font-weight: bold;
    cursor: pointer;
    transition: transform 0.1s;
}

.btn-plan:hover {
    background: #246343;
    transform: translateY(-2px);
}

.site-cell.clickable,
.token-cell.clickable {
    cursor: pointer;
    padding: 6px;
    border-radius: 6px;
    transition: background 0.2s;
    display: inline-block;
}

.site-cell.clickable:hover,
.token-cell.clickable:hover {
    background: #f0f7f3;
}

.site-cell small {
    color: #2d6a4f;
    font-weight: bold;
}

.history-bg {
    background: #fcfcfc;
}
.flex-between {
    display: flex;
    justify-content: space-between;
    align-items: center;
}
.search-mini input {
    border: 1px solid #eee;
    padding: 8px 14px;
    border-radius: 8px;
    font-size: 0.85rem;
    outline: none;
    width: 250px;
}
.search-mini input:focus {
    border-color: #2d7a4f;
}
.empty-text {
    color: #ccc;
    font-style: italic;
    font-size: 0.9rem;
    padding: 2rem 0;
    text-align: center;
}

.modal-overlay {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.4);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1000;
}
.site-modal {
    background: white;
    width: 90%;
    max-width: 420px;
    border-radius: 16px;
    padding: 2rem;
    box-shadow: 0 10px 30px rgba(0, 0, 0, 0.2);
}
.modal-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1.5rem;
}
.modal-header h3 {
    color: #2d6a4f;
    margin: 0;
    font-size: 1.2rem;
}
.close-btn {
    background: none;
    border: none;
    font-size: 1.8rem;
    cursor: pointer;
    color: #aaa;
}
.info-row {
    margin-bottom: 1.2rem;
}
.info-row label {
    display: block;
    font-size: 0.7rem;
    font-weight: 800;
    color: #aaa;
    margin-bottom: 4px;
}
.info-row p {
    margin: 0;
    color: #333;
    font-weight: 600;
    line-height: 1.5;
}
.modal-close-btn {
    width: 100%;
    margin-top: 1.5rem;
    padding: 12px;
}

.tag-don {
    background: #e8f5e9;
    color: #2e7d32;
    font-weight: 800;
    padding: 4px 10px;
    border-radius: 6px;
    font-size: 0.7rem;
}
.tag-vente {
    background: #fff8e1;
    color: #f57f17;
    font-weight: 800;
    padding: 4px 10px;
    border-radius: 6px;
    font-size: 0.7rem;
}

.small-date {
    font-family: monospace;
    color: #666;
    font-size: 0.95rem;
}

.btn-facture {
    background-color: #f0f4f1;
    color: #2d7a4f;
    border: 1px solid #2d7a4f;
    padding: 6px 12px;
    border-radius: 8px;
    font-size: 0.75rem;
    font-weight: 800;
    cursor: pointer;
    transition: all 0.2s;
    display: inline-flex;
    align-items: center;
    justify-content: center;
}

.btn-facture:hover {
    background-color: #2d7a4f;
    color: #ffffff;
    transform: translateY(-2px);
    box-shadow: 0 4px 8px rgba(45, 122, 79, 0.2);
}

.info-text {
    font-size: 0.9rem;
    color: #666;
    margin-bottom: 15px;
    line-height: 1.4;
}

.token-box {
    background: #f0f4f1;
    padding: 1.5rem;
    border-radius: 12px;
    font-family: monospace;
    font-size: 1.1rem;
    font-weight: bold;
    color: #1a1a1a;
    word-break: break-all;
    border: 2px dashed #2d6a4f;
    text-align: center;
    letter-spacing: 2px;
}
</style>
