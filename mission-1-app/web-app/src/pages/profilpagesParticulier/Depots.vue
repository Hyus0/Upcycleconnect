<template>
    <header class="content-header">
        <div class="header-left">
            <p class="sidebar__category2">ACCUEIL > DEPOTS</p>
            <h1 class="hero-title1">GESTION DES FLUX</h1>
            <p class="classic-text">
                Suivez vos objets, de la réservation à la récupération finale.
            </p>
        </div>
        <button
            class="btn-main-action"
            @click="$router.push('/profil/createAnnonce')"
        >
            + Déposer une annonce
        </button>
    </header>

    <div v-if="loading" class="loading-state">Chargement de vos flux...</div>

    <div v-else class="section-container">
        <div class="dash-block">
            <h2 class="block-title">
                ANNONCES À PLANIFIER
                <span class="badge">{{ aPlanifier.length }}</span>
            </h2>
            <table class="data-table">
                <thead>
                    <tr>
                        <th>OBJET</th>
                        <th>TYPE</th>
                        <th>VILLE ORIGINE</th>
                        <th class="text-right">ACTION</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-for="item in aPlanifier" :key="item.id">
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
                                        : "💰 " + item.prix + "€"
                                }}
                            </span>
                        </td>
                        <td>📍 {{ item.ville }}</td>
                        <td class="text-right">
                            <button
                                class="btn-plan"
                                @click="planifier(item.id)"
                            >
                                Réserver un casier
                            </button>
                        </td>
                    </tr>
                </tbody>
            </table>
            <p v-if="aPlanifier.length === 0" class="empty-text">
                Aucun objet en attente de réservation.
            </p>
        </div>

        <div class="dash-block">
            <h2 class="block-title">
                FLUX ACTIFS <span class="badge">{{ actifs.length }}</span>
            </h2>
            <table class="data-table table-active">
                <thead>
                    <tr>
                        <th>OBJET</th>
                        <th>SITE DE COLLECTE</th>
                        <th>CODE PIN</th>
                        <th>STATUT</th>
                        <th class="text-right">ACTION</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-for="item in actifs" :key="item.id">
                        <td>
                            <strong>{{ item.titre }}</strong>
                        </td>
                        <td>
                            <div
                                class="site-cell clickable"
                                @click="fetchAndOpenSite(item.id_site)"
                            >
                                <strong>{{
                                    sitesCache[item.id_site] || "Chargement..."
                                }}</strong
                                ><br />
                                <small>Voir les détails ↗</small>
                            </div>
                        </td>
                        <td>
                            <div
                                v-if="
                                    item.code_barre_retrait ||
                                    item.code_barre_depot
                                "
                                class="token-cell clickable"
                                @click.stop.prevent="
                                    showToken(
                                        item.code_barre_retrait
                                            ? item.code_barre_retrait
                                            : item.code_barre_depot,
                                    )
                                "
                            >
                                <span class="pin-badge">
                                    {{ truncateToken(item.code_barre_depot) }}
                                </span>
                            </div>
                            <span v-else>---</span>
                        </td>
                        <td>
                            <span
                                :class="[
                                    'status-pill',
                                    item.statut.toLowerCase(),
                                ]"
                            >
                                {{
                                    item.statut === "Reserve"
                                        ? "⌛ RÉSERVÉ"
                                        : "📦 DÉPOSÉ"
                                }}
                            </span>
                        </td>
                        <td class="text-right">
                            <button
                                class="btn-remove"
                                @click="annulerFlux(item.id, item.statut)"
                            >
                                {{
                                    item.statut === "Reserve"
                                        ? "Annuler"
                                        : "Retirer l'objet"
                                }}
                            </button>
                        </td>
                    </tr>
                </tbody>
            </table>
        </div>

        <div class="dash-block history-bg">
            <div class="flex-between">
                <h2 class="block-title">HISTORIQUE</h2>
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
                        <th>SITE</th>
                        <th class="text-center">STATUT FINAL</th>
                        <th class="text-right">DATE</th>
                        <th class="text-center">FACTURE</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-for="item in filteredHistory" :key="item.id">
                        <td>
                            <strong>{{ item.titre }}</strong>
                        </td>
                        <td>
                            <small>{{
                                sitesCache[item.id_site] || "Site inconnu"
                            }}</small>
                        </td>
                        <td class="text-center">
                            <span
                                :class="[
                                    'status-neutral',
                                    item.statut.toLowerCase(),
                                ]"
                            >
                                {{ item.statut }}
                            </span>
                        </td>
                        <td class="text-right small-date">
                            <span v-if="item.statut === 'Recupere'">
                                {{
                                    formatDate(item.date_recuperation_effective)
                                }}
                            </span>
                            <span v-else-if="item.statut === 'Depose'">
                                {{ formatDate(item.date_depot_effective) }}
                            </span>
                            <span v-else>
                                {{ formatDate(item.updated_at) }}
                            </span>
                        </td>
                        <td class="text-center">
                            <button
                                v-if="item.statut === 'Recupere'"
                                class="btn-facture"
                                @click="telechargerFacture(item.id)"
                                title="Télécharger la facture PDF"
                            >
                                Télécharger le PDF
                            </button>
                            <span v-else class="empty-cell">-</span>
                        </td>
                    </tr>
                </tbody>
            </table>
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
                        {{ selectedSite.adresse }}<br />{{
                            selectedSite.code_postal
                        }}
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
                <h3>Votre Token Sécurisé</h3>
                <button class="close-btn" @click="selectedToken = null">
                    ×
                </button>
            </div>
            <div class="modal-content">
                <p class="info-text">
                    Présentez ce code unique au niveau du conteneur sécurisé :
                </p>
                <div class="token-box">
                    {{ selectedToken }}
                </div>
            </div>
            <button
                class="btn-main-action modal-close-btn"
                @click="selectedToken = null"
            >
                Fermer
            </button>
        </div>
    </div>
</template>

<script setup>
import { ref, onMounted, computed } from "vue";
import { useRouter } from "vue-router";

const router = useRouter();
const annonces = ref([]);
const sitesCache = ref({});
const loading = ref(true);
const searchQuery = ref("");
const selectedSite = ref(null);
const selectedToken = ref(null);

const aPlanifier = computed(() =>
    annonces.value.filter((a) => a.statut === "Paye"),
);
const actifs = computed(() =>
    annonces.value.filter((a) => ["Reserve", "Depose"].includes(a.statut)),
);
const filteredHistory = computed(() => {
    const list = annonces.value.filter((a) =>
        ["Paye", "Recupere", "Annule"].includes(a.statut),
    );
    return list.filter((a) =>
        a.titre.toLowerCase().includes(searchQuery.value.toLowerCase()),
    );
});

const telechargerFacture = (idAnnonce) => {
    alert(
        `Génération de la facture pour l'objet #${idAnnonce} en cours...\n(À connecter avec le backend Go 🚀)`,
    );
};

const formatDate = (d) => {
    if (!d || d.startsWith("0001")) return "...";
    return new Date(d).toLocaleDateString("fr-FR", {
        day: "numeric",
        month: "short",
        year: "numeric",
    });
};

const fetchSitesCache = async () => {
    try {
        const res = await fetch("http://localhost:8081/sites");
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

const fetchAnnonces = async () => {
    const id = sessionStorage.getItem("userId");
    const token = sessionStorage.getItem("userToken");
    loading.value = true;
    try {
        const res = await fetch(`http://localhost:8081/users/${id}/annonces`, {
            headers: { Authorization: token },
        });
        if (res.ok) annonces.value = await res.json();
    } finally {
        loading.value = false;
    }
};

const fetchAndOpenSite = async (siteId) => {
    if (!siteId) return;
    const token = sessionStorage.getItem("userToken");
    try {
        const res = await fetch(`http://localhost:8081/site/${siteId}`, {
            headers: { Authorization: token },
        });
        if (res.ok) {
            selectedSite.value = await res.json();
        }
    } catch (e) {
        console.error(e);
    }
};

const planifier = (id) =>
    router.push({ name: "reserve-casier", params: { id } });

const annulerFlux = async (id, statut) => {
    const actionLabel =
        statut === "Reserve" ? "annuler la réservation" : "retirer l'objet";

    if (!confirm(`Voulez-vous vraiment ${actionLabel} ?`)) return;

    const token = sessionStorage.getItem("userToken");

    try {
        const res = await fetch(
            `http://localhost:8081/annonces/${id}/retirer`,
            {
                method: "POST",
                headers: {
                    Authorization: token,
                    "Content-Type": "application/json",
                },
            },
        );

        if (res.ok) {
            alert(
                statut === "Reserve"
                    ? "Réservation annulée et casier libéré."
                    : "Objet retiré. L'annonce est de nouveau disponible !",
            );

            fetchAnnonces();
        } else {
            const error = await res.text();
            alert("Erreur lors de l'opération : " + error);
        }
    } catch (e) {
        console.error("Erreur serveur :", e);
        alert("Impossible de joindre le serveur.");
    }
};

onMounted(() => {
    fetchSitesCache();
    fetchAnnonces();
});
</script>

<style scoped>
.dash-block {
    background: white;
    border: 1px solid #eee;
    border-radius: 16px;
    padding: 1.5rem;
    margin-bottom: 2rem;
}

.block-title {
    font-size: 0.75rem;
    font-weight: 800;
    color: #999;
    margin-bottom: 1.5rem;
    display: flex;
    align-items: center;
    gap: 10px;
}

.badge {
    background: #f0f0f0;
    color: #666;
    padding: 2px 8px;
    border-radius: 10px;
    font-size: 0.7rem;
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
    text-align: right;
}

.pin-badge {
    background: #e8f5e9;
    color: #2d6a4f;
    font-weight: 900;
    padding: 4px 10px;
    border-radius: 6px;
    letter-spacing: 1px;
}

.status-pill {
    font-size: 0.7rem;
    font-weight: bold;
    padding: 4px 8px;
    border-radius: 4px;
}
.status-pill.reserve {
    background: #fff3e0;
    color: #e65100;
}
.status-pill.depose {
    background: #e3f2fd;
    color: #1976d2;
}

.status-neutral {
    font-size: 0.75rem;
    font-weight: 700;
    text-transform: uppercase;
}
.status-neutral.recupere {
    color: #2e7d32;
}
.status-neutral.paye {
    color: #f57f17;
}
.status-neutral.annule {
    color: #dc3545;
}

.btn-plan {
    background: #2d6a4f;
    color: white;
    border: none;
    padding: 7px 15px;
    border-radius: 8px;
    font-weight: bold;
    cursor: pointer;
}
.btn-remove {
    background: none;
    border: none;
    color: #dc3545;
    font-weight: bold;
    cursor: pointer;
    font-size: 0.8rem;
}

.site-cell.clickable,
.token-cell.clickable {
    cursor: pointer;
    padding: 4px;
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
    padding: 6px 12px;
    border-radius: 8px;
    font-size: 0.8rem;
    outline: none;
}
.empty-text {
    color: #ccc;
    font-style: italic;
    font-size: 0.85rem;
    padding: 1rem 0;
}

.modal-overlay {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.3);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1000;
}
.site-modal {
    background: white;
    width: 90%;
    max-width: 400px;
    border-radius: 20px;
    padding: 2rem;
    box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
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
}
.close-btn {
    background: none;
    border: none;
    font-size: 1.5rem;
    cursor: pointer;
    color: #ccc;
}
.info-row {
    margin-bottom: 1.2rem;
}
.info-row label {
    display: block;
    font-size: 0.65rem;
    font-weight: 800;
    color: #aaa;
    margin-bottom: 4px;
}
.info-row p {
    margin: 0;
    color: #333;
    font-weight: 600;
    line-height: 1.4;
}
.modal-close-btn {
    width: 100%;
    margin-top: 1rem;
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

.text-center {
    text-align: center !important;
}
.text-right {
    text-align: right !important;
}
.small-date {
    font-family: monospace;
    color: #666;
    font-size: 0.85rem;
}
.empty-cell {
    color: #ccc;
    font-weight: bold;
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
    transform: translateY(-1px);
    box-shadow: 0 4px 8px rgba(45, 122, 79, 0.2);
}

.info-text {
    font-size: 0.85rem;
    color: #666;
    margin-bottom: 15px;
}

.token-box {
    background: #f0f4f1;
    padding: 1.5rem;
    border-radius: 12px;
    font-family: monospace;
    font-size: 0.95rem;
    color: #1a1a1a;
    word-break: break-all;
    border: 1px dashed #2d6a4f;
    text-align: center;
    letter-spacing: 1px;
    line-height: 1.5;
}
</style>
