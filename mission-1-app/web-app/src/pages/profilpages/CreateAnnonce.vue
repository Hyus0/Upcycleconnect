<template>
    <header class="content-header">
        <div class="header-left">
            <p class="sidebar__category2">ACCUEIL &gt; ANNONCES &gt; CREER</p>
            <h1 class="hero-title1">DEPOSER UN OBJET</h1>
            <p class="classic-text">Cette annonce sera creee dans la base et visible via les parcours front pilotes par le back.</p>
        </div>
        <button class="btn-secondary" @click="$router.back()">Retour</button>
    </header>

    <div v-if="errors.length" class="error-box"><ul><li v-for="(err, index) in errors" :key="index">{{ err }}</li></ul></div>
    <div v-if="successMsg" class="success-box">{{ successMsg }}</div>

    <form class="create-annonce-form" @submit.prevent="handleSubmit">
        <div class="split-layout">
            <div class="form-card">
                <h2 class="card-title">1. L'objet</h2>
                <div class="form-group"><label>Titre</label><input v-model="form.titre" type="text" required /></div>
                <div class="form-group"><label>Description</label><textarea v-model="form.description" rows="6"></textarea></div>
                <div class="form-row">
                    <div class="form-group">
                        <label>Categorie</label>
                        <select v-model="form.id_categorie" required>
                            <option value="" disabled>Choisir une categorie</option>
                            <option v-for="cat in categories" :key="cat.id" :value="cat.id">{{ cat.nom }}</option>
                        </select>
                    </div>
                    <div class="form-group">
                        <label>Etat</label>
                        <select v-model="form.etat_objet">
                            <option value="Neuf">Neuf</option>
                            <option value="Bon etat">Bon etat</option>
                            <option value="Usage">Usage</option>
                        </select>
                    </div>
                </div>
                <div class="form-row">
                    <div class="form-group">
                        <label>Type</label>
                        <select v-model="form.type">
                            <option value="Don">Don</option>
                            <option value="Vente">Vente</option>
                        </select>
                    </div>
                    <div class="form-group" v-if="form.type === 'Vente'">
                        <label>Prix</label>
                        <input v-model.number="form.prix" type="number" min="0" step="0.01" />
                    </div>
                </div>
            </div>

            <div class="form-card">
                <h2 class="card-title">2. Details et localisation</h2>
                <div class="form-row">
                    <div class="form-group"><label>Materiau</label><input v-model="form.type_materiau" type="text" /></div>
                    <div class="form-group"><label>Poids (kg)</label><input v-model.number="form.poids_estime_kg" type="number" step="0.1" min="0" /></div>
                </div>
                <div class="form-row">
                    <div class="form-group"><label>Ville</label><input v-model="form.ville" type="text" required /></div>
                    <div class="form-group"><label>Code postal</label><input v-model="form.code_postal" type="text" required /></div>
                </div>
                <div class="form-group"><label>Adresse</label><input v-model="form.adresse" type="text" /></div>
                <div class="form-actions-card">
                    <button type="button" class="btn-cancel" @click="$router.back()">Annuler</button>
                    <button type="submit" class="btn-save" :disabled="loading || !categories.length">{{ loading ? "Publication..." : "Publier l'annonce" }}</button>
                </div>
            </div>
        </div>
    </form>
</template>

<script setup>
import { onMounted, ref, watch } from "vue";
import { useRouter } from "vue-router";
import { createAnnonce, fetchCategories } from "../../services/publicApi";

const router = useRouter();
const loading = ref(false);
const categories = ref([]);
const errors = ref([]);
const successMsg = ref("");

const form = ref({
    id_vendeur: Number(localStorage.getItem("userId")) || 0,
    id_categorie: "",
    titre: "",
    description: "",
    type_materiau: "",
    poids_estime_kg: 0,
    prix: 0,
    etat_objet: "Bon etat",
    type: "Don",
    ville: "",
    code_postal: "",
    adresse: "",
    statut: "Disponible"
});

watch(
    () => form.value.type,
    (nextType) => {
        if (nextType === "Don") {
            form.value.prix = 0;
        }
    }
);

const loadCategories = async () => {
    try {
        const payload = await fetchCategories();
        categories.value = Array.isArray(payload) ? payload : [];
        if (!categories.value.length) {
            errors.value = ["Aucune categorie n'est disponible pour le moment."];
        }
    } catch (error) {
        console.error("Erreur categories :", error);
        categories.value = [];
        errors.value = [error.message || "Impossible de charger les categories."];
    }
};

const validateFrontend = () => {
    errors.value = [];
    form.value.id_vendeur = Number(localStorage.getItem("userId")) || 0;

    if (!form.value.id_vendeur) errors.value.push("Vous devez etre connecte pour deposer une annonce.");
    if ((form.value.titre || "").length < 5) errors.value.push("Le titre doit faire au moins 5 caracteres.");
    if ((form.value.description || "").length < 10) errors.value.push("La description doit faire au moins 10 caracteres.");
    if (!form.value.id_categorie) errors.value.push("Veuillez choisir une categorie.");
    if (form.value.type === "Vente" && Number(form.value.prix) <= 0) errors.value.push("Le prix doit etre superieur a 0 EUR.");
    if ((form.value.code_postal || "").length !== 5) errors.value.push("Le code postal doit contenir 5 caracteres.");
    if ((form.value.ville || "").trim().length < 2) errors.value.push("La ville doit etre renseignee.");
    return errors.value.length === 0;
};

const buildPayload = () => ({
    ...form.value,
    id_vendeur: Number(localStorage.getItem("userId")) || 0,
    id_categorie: Number(form.value.id_categorie) || 0,
    prix: form.value.type === "Don" ? 0 : Number(form.value.prix) || 0,
    poids_estime_kg: Number(form.value.poids_estime_kg) || 0,
    titre: (form.value.titre || "").trim(),
    description: (form.value.description || "").trim(),
    type_materiau: (form.value.type_materiau || "").trim(),
    ville: (form.value.ville || "").trim(),
    code_postal: (form.value.code_postal || "").trim(),
    adresse: (form.value.adresse || "").trim(),
    statut: "Disponible"
});

const handleSubmit = async () => {
    if (!validateFrontend()) return;
    loading.value = true;
    successMsg.value = "";
    try {
        await createAnnonce(buildPayload());
        successMsg.value = "Annonce creee avec succes.";
        setTimeout(() => router.push("/profil/annonces"), 1000);
    } catch (error) {
        errors.value = [error.message || "La creation a echoue."];
    } finally {
        loading.value = false;
    }
};

onMounted(loadCategories);
</script>

<style scoped>
.error-box { background-color: #fee2e2; border: 1px solid #ef4444; color: #b91c1c; padding: 12px; border-radius: 10px; margin-bottom: 15px; }
.success-box { background-color: #f0fdf4; border: 1px solid #22c55e; color: #166534; padding: 12px; border-radius: 10px; margin-bottom: 15px; }
.split-layout { display: grid; grid-template-columns: 1.5fr 1fr; gap: 2rem; width: 100%; }
.form-card { background: #ffffff; padding: 2rem; border-radius: 16px; border: 1px solid #eee; display: flex; flex-direction: column; gap: 1.5rem; margin-bottom: 2rem; }
.form-row { display: grid; grid-template-columns: 1fr 1fr; gap: 1.5rem; }
.form-group { display: flex; flex-direction: column; gap: 0.6rem; }
input, select, textarea { padding: 0.9rem; border: 1px solid #ddd; border-radius: 10px; font-family: inherit; font-size: 0.95rem; }
.form-actions-card { display: flex; justify-content: flex-end; gap: 1rem; margin-top: auto; }
</style>
