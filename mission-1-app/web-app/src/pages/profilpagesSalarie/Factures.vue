<template>
    <main class="page-container">
        <header class="content-header">
            <div class="header-left">
                <p class="sidebar__category2">
                    ACCUEIL > MON COMPTE > ESPACE PRO
                </p>
                <h1 class="hero-title1">Mes Factures</h1>
                <p class="classic-text">
                    Retrouvez ici l'historique complet de vos transactions et
                    téléchargez vos factures.
                </p>
            </div>
            <div class="header-actions">
                <button class="btn-view" type="button" @click="fetchFactures">
                    Actualiser
                </button>
                <button
                    class="btn-secondary"
                    type="button"
                    @click="$router.push('/profil')"
                >
                    Retour au profil
                </button>
            </div>
        </header>
        <div v-if="loading" class="state-card">
            Chargement de vos factures...
        </div>
        <div v-else-if="error" class="state-card error">{{ error }}</div>
        <div v-else-if="factures.length === 0" class="empty-state">
            <div class="empty-icon">📄</div>
            <h2>Aucune facture disponible</h2>
            <p>
                Vos factures seront générées automatiquement après chaque
                commande validée.
            </p>
        </div>

        <div v-else class="factures-list">
            <article
                v-for="facture in factures"
                :key="facture.id"
                class="facture-row"
            >
                <div class="facture-info">
                    <span class="invoice-number">{{
                        facture.numero_facture
                    }}</span>
                    <h3 class="facture-title">
                        Commande #{{ facture.commande_id }}
                    </h3>
                    <p class="facture-meta">
                        {{ facture.date_transaction }} •
                        {{ facture.statut_paiement }}
                    </p>
                </div>

                <div class="facture-amount">
                    <strong>{{ formatPrice(facture.montant_total) }}</strong>
                </div>

                <div class="row-actions">
                    <button
                        class="btn-view"
                        type="button"
                        @click="downloadFacture(facture.id)"
                    >
                        Télécharger
                    </button>
                    <button
                        class="btn-secondary"
                        type="button"
                        @click="sendFacture(facture.id)"
                    >
                        Envoyer par mail
                    </button>
                </div>
            </article>
        </div>
    </main>
</template>

<script setup>
import { onMounted, ref } from "vue";

const API_URL = "http://localhost:8081";
const loading = ref(true);
const error = ref("");
const factures = ref([]);

const userId = () =>
    Number(sessionStorage.getItem("id") || sessionStorage.getItem("userId")) ||0;

const formatPrice = (value) =>
    new Intl.NumberFormat("fr-FR", {
        style: "currency",
        currency: "EUR",
    }).format(Number(value) || 0);

const fetchFactures = async () => {
    loading.value = true;
    error.value = "";
    try {
        const res = await fetch(`${API_URL}/users/${userId()}/factures`, {
            headers: {
                Authorization: sessionStorage.getItem("userToken") || "",
            },
        });
        if (!res.ok) throw new Error("Impossible de charger les factures.");
        factures.value = (await res.json()) || [];
    } catch (err) {
        error.value = err.message;
    } finally {
        loading.value = false;
    }
};

const downloadFacture = (factureId) => {
    window.open(
        `${API_URL}/users/${userId()}/factures/${factureId}/download`,
        "_blank",
    );
};

const sendFacture = async (factureId) => {
    try {
        const res = await fetch(
            `${API_URL}/users/${userId()}/factures/${factureId}/send`,
            {
                method: "POST",
                headers: {
                    Authorization: sessionStorage.getItem("userToken") || "",
                },
            },
        );
        const payload = res.ok
            ? await res.json()
            : { message: "Erreur lors de l'envoi." };
        alert(payload.message);
    } catch (error) {
        alert("Impossible de contacter le serveur.");
    }
};

onMounted(fetchFactures);
</script>

<style scoped>
.page-main-content {
    padding: 40px;
    max-width: 1100px;
    margin: 0 auto;
    font-family: "Syne", sans-serif;
}

.content-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-end;
    margin-bottom: 3rem;
}

.sidebar__category2 {
    padding-top: 20px;
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
    color: #6d7b72;
    margin-top: 10px;
}

.factures-list {
    display: flex;
    flex-direction: column;
    gap: 16px;
}

.facture-row {
    background: #fff;
    border: 1px solid #e5ede7;
    border-radius: 16px;
    padding: 24px;
    display: grid;
    grid-template-columns: 1fr auto auto;
    align-items: center;
    gap: 20px;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.02);
}

.invoice-number {
    color: #2f8f5b;
    font-size: 0.75rem;
    font-weight: 900;
    letter-spacing: 1px;
    text-transform: uppercase;
}

.facture-title {
    margin: 4px 0;
    font-size: 1.1rem;
}
.facture-meta {
    color: #8fa396;
    font-size: 0.85rem;
    margin: 0;
}
.facture-amount strong {
    font-size: 1.2rem;
    color: #2d7a4f;
}

.row-actions {
    display: flex;
    gap: 10px;
}

button {
    border-radius: 10px;
    padding: 10px 16px;
    font-weight: 800;
    cursor: pointer;
    border: none;
    font-family: inherit;
    transition: 0.2s;
}

.btn-view {
    background: #2f8f5b;
    color: #fff;
}
.btn-secondary {
    background: #f0f4f1;
    color: #2d7a4f;
}
.btn-view:hover {
    background: #23653e;
}
.btn-secondary:hover {
    background: #e1ede5;
}

.state-card,
.empty-state {
    padding: 40px;
    text-align: center;
    border: 1px dashed #cfe0d4;
    border-radius: 16px;
    color: #63746a;
}

@media (max-width: 760px) {
    .facture-row {
        grid-template-columns: 1fr;
        text-align: center;
    }
    .row-actions {
        justify-content: center;
    }
}
</style>
