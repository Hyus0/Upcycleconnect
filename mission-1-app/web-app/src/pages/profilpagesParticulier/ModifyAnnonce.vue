<template>
    <header class="content-header">
        <div class="header-left">
            <p class="sidebar__category2">ACCUEIL > MES ANNONCES > ÉDITION</p>
            <h1 class="hero-title1">MODIFIER L'OBJET</h1>
        </div>
        <div class="header-actions">
            <button class="btn-secondary" @click="$router.back()">
                🠔 Retour
            </button>
            <button
                type="submit"
                form="editForm"
                class="btn-main-action-header"
                :disabled="submitting"
            >
                {{ submitting ? "..." : "Enregistrer" }}
            </button>
        </div>
    </header>

    <div v-if="loading" class="loading-state">Récupération des données...</div>

    <form
        v-else
        id="editForm"
        @submit.prevent="handleUpdate"
        class="split-layout"
    >
        <div class="left-column">
            <div class="info-card">
                <div class="card-header-flex">
                    <h2 class="card-title">Détails de l'annonce</h2>
                    <span
                        :class="
                            annonce.type === 'Don' ? 'tag-don' : 'tag-vente'
                        "
                        class="type-badge"
                    >
                        {{ annonce.type === "Don" ? "🎁 DON" : "💰 VENTE" }}
                    </span>
                </div>

                <div class="form-group">
                    <label class="info-label">Titre de l'annonce</label>
                    <input
                        v-model="annonce.titre"
                        type="text"
                        class="input-field"
                        required
                    />
                </div>

                <div class="form-group">
                    <label class="info-label">Description détaillée</label>
                    <textarea
                        v-model="annonce.description"
                        class="edit-textarea"
                        rows="8"
                    ></textarea>
                </div>

                <div class="specs-grid-edit">
                    <div class="form-group">
                        <label class="info-label">Catégorie</label>
                        <select
                            v-model="annonce.id_categorie"
                            class="input-field"
                        >
                            <option
                                v-for="cat in categories"
                                :key="cat.id"
                                :value="cat.id"
                            >
                                {{ cat.nom }}
                            </option>
                        </select>
                    </div>
                    <div class="form-group">
                        <label class="info-label">État de l'objet</label>
                        <select
                            v-model="annonce.etat_objet"
                            class="input-field"
                        >
                            <option value="Neuf">Neuf</option>
                            <option value="Bon etat">Bon état</option>
                            <option value="Usage">Usagé</option>
                        </select>
                    </div>
                </div>

                <div class="specs-grid-edit">
                    <div class="form-group">
                        <label class="info-label">Matériau</label>
                        <input
                            v-model="annonce.type_materiau"
                            type="text"
                            class="input-field"
                            placeholder="Ex: Bois, Acier..."
                        />
                    </div>
                    <div class="form-group">
                        <label class="info-label">Poids estimé (kg)</label>
                        <input
                            v-model.number="annonce.poids_estime_kg"
                            type="number"
                            step="0.1"
                            class="input-field"
                        />
                    </div>
                </div>
            </div>
        </div>

        <div class="right-column">
            <div class="status-card" v-if="annonce.type === 'Vente'">
                <h3>Transaction</h3>
                <div class="form-group">
                    <label class="info-label">Prix de vente (€)</label>
                    <input
                        v-model.number="annonce.prix"
                        type="number"
                        step="0.01"
                        class="input-field"
                    />
                </div>
            </div>

            <div class="info-card">
                <h3>📍 Localisation</h3>
                <div class="form-group">
                    <label class="info-label">Ville</label>
                    <input
                        v-model="annonce.ville"
                        type="text"
                        class="input-field"
                    />
                </div>

                <div class="specs-grid-edit">
                    <div class="form-group">
                        <label class="info-label">Code Postal</label>
                        <input
                            v-model="annonce.code_postal"
                            type="text"
                            class="input-field"
                            maxlength="5"
                        />
                    </div>
                </div>

                <div class="form-group">
                    <label class="info-label">Adresse complète</label>
                    <input
                        v-model="annonce.adresse"
                        type="text"
                        class="input-field"
                    />
                </div>
            </div>
        </div>
    </form>
</template>

<script setup>
import { ref, onMounted } from "vue";
import { useRoute, useRouter } from "vue-router";

const route = useRoute();
const router = useRouter();
const loading = ref(true);
const submitting = ref(false);
const categories = ref([]);
const annonce = ref({});

const fetchCategories = async () => {
    try {
        const res = await fetch("http://localhost:8081/categories");
        if (res.ok) categories.value = await res.json();
    } catch (e) {
        console.error("Erreur catégories", e);
    }
};

const fetchAnnonce = async () => {
    const token = sessionStorage.getItem("userToken");
    try {
        const res = await fetch(
            `http://localhost:8081/annonces/${route.params.id}`,
            {
                headers: { Authorization: token },
            },
        );
        if (res.ok) annonce.value = await res.json();
    } finally {
        loading.value = false;
    }
};

onMounted(async () => {
    await Promise.all([fetchCategories(), fetchAnnonce()]);
});

const handleUpdate = async () => {
    submitting.value = true;
    const token = sessionStorage.getItem("userToken");
    try {
        const res = await fetch(
            `http://localhost:8081/annonces/${annonce.value.id}`,
            {
                method: "PUT",
                headers: {
                    Authorization: token,
                    "Content-Type": "application/json",
                },
                body: JSON.stringify(annonce.value),
            },
        );
        if (res.ok)
            router.push({
                name: "see-annonce",
                params: { id: annonce.value.id },
            });
    } finally {
        submitting.value = false;
    }
};
</script>

<style scoped>
.split-layout {
    display: grid;
    grid-template-columns: 1.5fr 1fr;
    gap: 2rem;
}

.info-card,
.status-card {
    background: white;
    padding: 1.5rem;
    border-radius: 16px;
    border: 1px solid #eee;
    margin-bottom: 1.5rem;
}

.card-header-flex {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 2rem;
    border-bottom: 1px solid #f5f5f5;
    padding-bottom: 1rem;
}

.type-badge {
    padding: 6px 14px;
    border-radius: 20px;
    font-weight: 800;
    font-size: 0.85rem;
}
.tag-don {
    background: #e8f5e9;
    color: #2e7d32;
}
.tag-vente {
    background: #fff8e1;
    color: #f57f17;
}

.form-group {
    margin-bottom: 1.5rem;
}

.info-label {
    display: block;
    font-size: 0.75rem;
    font-weight: 700;
    text-transform: uppercase;
    color: #aaa;
    margin-bottom: 8px;
    letter-spacing: 0.5px;
}

.input-field {
    width: 100%;
    padding: 0.85rem;
    border-radius: 10px;
    border: 1px solid #ddd;
    background: #fcfcfc;
    font-family: inherit;
    font-size: 1rem;
    transition: all 0.2s;
}

.input-field:focus {
    border-color: #2d6a4f;
    outline: none;
    background: white;
    box-shadow: 0 0 0 3px rgba(45, 106, 79, 0.05);
}

.edit-textarea {
    width: 100%;
    padding: 1rem;
    border-radius: 12px;
    border: 1px solid #eee;
    background: #f9f9f9;
    resize: vertical;
    line-height: 1.6;
    font-family: inherit;
}

.specs-grid-edit {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 1.5rem;
}

.content-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 2rem;
}
.header-actions {
    display: flex;
    gap: 1rem;
}

.btn-main-action-header {
    background: #2d6a4f;
    color: white;
    padding: 0.6rem 1.5rem;
    border-radius: 10px;
    font-weight: bold;
    cursor: pointer;
    border: none;
    transition: background 0.2s;
}
.btn-main-action-header:hover {
    background: #1b4332;
}

h3 {
    font-size: 0.85rem;
    text-transform: uppercase;
    color: #999;
    margin-bottom: 1.2rem;
}
</style>
