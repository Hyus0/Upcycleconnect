<template>
    <header class="content-header">
        <div class="header-left">
            <p class="sidebar__category2">
                ACCUEIL > MES ANNONCES > {{ annonce?.titre }}
            </p>
            <h1 class="hero-title1">{{ annonce?.titre || "Chargement..." }}</h1>
        </div>
        <div class="header-actions">
            <button class="btn-secondary" @click="$router.back()">
                🠔 Retour
            </button>
            <button
                v-if="annonce?.est_valide === 'En attente'"
                class="btn-modify"
                @click="goToEdit"
            >
                Modifier l'annonce
            </button>
        </div>
    </header>

    <div v-if="loading" class="loading-state">Récupération des données...</div>

    <div v-else-if="annonce?.id" class="split-layout">
        <div class="left-column">
            <div class="info-card">
                <div class="card-header-flex">
                    <h2 class="card-title">Détails de l'objet</h2>
                    <span
                        :class="
                            annonce.type === 'Don' ? 'tag-don' : 'tag-vente'
                        "
                        class="type-badge"
                    >
                        {{
                            annonce.type === "Don"
                                ? "🎁 DON"
                                : "💰 VENTE " + annonce.prix + "€"
                        }}
                    </span>
                </div>

                <div class="description-section">
                    <label class="info-label">Description</label>
                    <div class="description-box">
                        {{ annonce.description }}
                    </div>
                </div>

                <div class="specs-grid">
                    <div class="spec-item">
                        <label>Catégorie</label>
                        <p class="highlight-val">
                            {{ categoryName || "Chargement..." }}
                        </p>
                    </div>
                    <div class="spec-item">
                        <label>Matériau</label>
                        <p>{{ annonce.type_materiau || "Non renseigné" }}</p>
                    </div>
                    <div class="spec-item">
                        <label>Poids estimé</label>
                        <p>{{ annonce.poids_estime_kg }} kg</p>
                    </div>
                    <div class="spec-item">
                        <label>État déclaré</label>
                        <p>{{ annonce.etat_objet }}</p>
                    </div>
                </div>

                <div class="location-section">
                    <label>Point de dépôt prévu</label>
                    <p class="address-text">
                        <strong>📍 {{ annonce.adresse }}</strong
                        ><br />
                        {{ annonce.code_postal }} {{ annonce.ville }}
                    </p>
                </div>
            </div>
        </div>

        <div class="right-column">
            <div class="status-card">
                <h3>Position & Validation</h3>
                <div class="data-row">
                    <span class="data-label">Vérification Check :</span>
                    <span
                        :class="
                            annonce.est_valide === 'Valide'
                                ? 'text-success'
                                : 'text-pending'
                        "
                    >
                        {{
                            annonce.est_valide === "Valide"
                                ? "✓ APPROUVÉ"
                                : "⌛ EN ANALYSE"
                        }}
                    </span>
                </div>
                <div class="data-row">
                    <span class="data-label">État actuel :</span>
                    <span class="status-badge">{{ annonce.statut }}</span>
                </div>
                <div v-if="annonce.id_casier" class="data-row">
                    <span class="data-label">Casier réservé :</span>
                    <span class="value-casier">N°{{ annonce.id_casier }}</span>
                </div>
            </div>

            <div class="dates-card">
                <h3>Historique logistique</h3>
                <div class="timeline-item">
                    <label>Mise en ligne</label>
                    <p>{{ formatDate(annonce.date_creation) }}</p>
                </div>
                <div class="timeline-item">
                    <label>Dépôt au conteneur</label>
                    <p
                        :class="
                            !annonce.date_depot_effective ? 'not-done' : 'done'
                        "
                    >
                        {{
                            annonce.date_depot_effective
                                ? formatDate(annonce.date_depot_effective)
                                : "❌ Pas encore déposé"
                        }}
                    </p>
                </div>
                <div class="timeline-item">
                    <label>Retrait par l'artisan</label>
                    <p
                        :class="
                            !annonce.date_recuperation_effective
                                ? 'not-done'
                                : 'done'
                        "
                    >
                        {{
                            annonce.date_recuperation_effective
                                ? formatDate(
                                      annonce.date_recuperation_effective,
                                  )
                                : "⌛ En attente de retrait"
                        }}
                    </p>
                </div>
            </div>

            <div
                v-if="annonce.code_pin_depot && annonce.statut !== 'Depose'"
                class="pin-card"
            >
                <label>CODE PIN D'OUVERTURE</label>
                <div class="pin-code">{{ annonce.code_pin_depot }}</div>
                <p class="pin-hint">
                    Valable au conteneur de {{ annonce.ville }}
                </p>
            </div>

            <button
                v-if="
                    annonce.est_valide === 'Valide' && !annonce.code_pin_depot
                "
                class="btn-main-action"
                @click="handleRequestDeposit"
            >
                📦 Demander l'accès au conteneur
            </button>

            <div v-if="annonce.id_acheteur" class="buyer-card">
                <p><strong>Objet réservé !</strong></p>
                <small>Artisan repreneur : #{{ annonce.id_acheteur }}</small>
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref, onMounted } from "vue";
import { useRoute, useRouter } from "vue-router";

const route = useRoute();
const router = useRouter();
const loading = ref(true);
const annonce = ref(null);
const categoryName = ref("");

const formatDate = (d) => {
    if (!d || d.startsWith("0001")) return null;
    return new Date(d).toLocaleString("fr-FR", {
        day: "numeric",
        month: "short",
        year: "numeric",
        hour: "2-digit",
        minute: "2-digit",
    });
};

const goToEdit = () => {
    router.push({
        name: "modification-annonce",
        params: { id: annonce.value.id },
    });
};

const fetchCategory = async (catId) => {
    try {
        const res = await fetch(`http://localhost:8081/category/${catId}`);
        if (res.ok) {
            const catData = await res.json();
            categoryName.value = catData.nom;
        }
    } catch (e) {
        console.error(e);
    }
};

const fetchAnnonce = async () => {
    const token = localStorage.getItem("userToken");
    try {
        const res = await fetch(
            `http://localhost:8081/annonces/${route.params.id}`,
            {
                headers: { Authorization: token },
            },
        );
        if (res.ok) {
            const data = await res.json();
            annonce.value = data;
            if (data.id_categorie) fetchCategory(data.id_categorie);
        }
    } finally {
        loading.value = false;
    }
};

const handleRequestDeposit = async () => {
    fetchAnnonce(); 
};

onMounted(fetchAnnonce);
</script>

<style scoped>
.card-header-flex {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1.5rem;
    border-bottom: 1px solid #f5f5f5;
    padding-bottom: 1rem;
}

.type-badge {
    padding: 6px 14px;
    border-radius: 20px;
    font-weight: 800;
    font-size: 0.85rem;
}

.description-section {
    margin-bottom: 2rem;
}

.info-label {
    display: block;
    font-size: 0.75rem;
    font-weight: 700;
    text-transform: uppercase;
    color: #999;
    margin-bottom: 8px;
    letter-spacing: 0.5px;
}

.description-box {
    background: #f9f9f9;
    padding: 1.2rem;
    border-radius: 12px;
    color: #444;
    line-height: 1.6;
    font-size: 1rem;
    border: 1px solid #f0f0f0;
}

.specs-grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 1.5rem;
    background: #fcfcfc;
    padding: 1.5rem;
    border-radius: 12px;
    border: 1px solid #f5f5f5;
    margin-bottom: 2rem;
}

.spec-item label {
    font-size: 0.7rem;
    color: #aaa;
    text-transform: uppercase;
    margin-bottom: 4px;
    display: block;
}

.spec-item p {
    font-weight: bold;
    color: #333;
    margin: 0;
}

.highlight-val {
    color: #2d6a4f;
}

.split-layout {
    display: grid;
    grid-template-columns: 1.5fr 1fr;
    gap: 2rem;
}
.info-card,
.status-card,
.dates-card,
.pin-card {
    background: white;
    padding: 1.5rem;
    border-radius: 16px;
    border: 1px solid #eee;
    margin-bottom: 1.5rem;
}
.tag-don {
    background: #e8f5e9;
    color: #2e7d32;
}
.tag-vente {
    background: #fff8e1;
    color: #f57f17;
}
.text-success {
    color: #2d6a4f;
    font-weight: bold;
}
.text-pending {
    color: #f57c00;
    font-weight: bold;
}
.timeline-item {
    margin-bottom: 1.2rem;
    padding-left: 15px;
    border-left: 2px solid #f0f0f0;
}
.pin-card {
    background: #1a1c23;
    color: white;
    text-align: center;
    border: none;
}
.pin-code {
    font-size: 3.5rem;
    font-weight: 900;
    color: #ffd43b;
    letter-spacing: 8px;
    margin: 10px 0;
}
.btn-main-action {
    width: 100%;
    background: #2d6a4f;
    color: white;
    padding: 1.2rem;
    border-radius: 12px;
    font-weight: bold;
    cursor: pointer;
    border: none;
}
</style>
