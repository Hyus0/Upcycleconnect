<template>
    <header class="content-header">
        <div class="header-left">
            <p class="sidebar__category2">ACCUEIL &gt; MES ANNONCES &gt; EDITION</p>
            <h1 class="hero-title1">MODIFIER L'OBJET</h1>
        </div>
        <div class="header-actions">
            <button class="btn-secondary" @click="$router.back()">Retour</button>
            <button type="submit" form="editForm" class="btn-main-action-header" :disabled="submitting">
                {{ submitting ? "..." : "Enregistrer" }}
            </button>
        </div>
    </header>

    <div v-if="loading" class="loading-state">Recuperation des donnees...</div>

    <form v-else id="editForm" class="split-layout" @submit.prevent="handleUpdate">
        <div class="info-card">
            <div class="form-group"><label class="info-label">Titre</label><input v-model="annonce.titre" type="text" class="input-field" required /></div>
            <div class="form-group"><label class="info-label">Description</label><textarea v-model="annonce.description" class="edit-textarea" rows="8"></textarea></div>
            <div class="specs-grid-edit">
                <div class="form-group">
                    <label class="info-label">Categorie</label>
                    <select v-model="annonce.id_categorie" class="input-field">
                        <option v-for="cat in categories" :key="cat.id" :value="cat.id">{{ cat.nom }}</option>
                    </select>
                </div>
                <div class="form-group">
                    <label class="info-label">Etat</label>
                    <select v-model="annonce.etat_objet" class="input-field">
                        <option value="Neuf">Neuf</option>
                        <option value="Bon etat">Bon etat</option>
                        <option value="Usage">Usage</option>
                    </select>
                </div>
            </div>
        </div>

        <div class="info-card">
            <div class="specs-grid-edit">
                <div class="form-group"><label class="info-label">Materiau</label><input v-model="annonce.type_materiau" type="text" class="input-field" /></div>
                <div class="form-group"><label class="info-label">Poids</label><input v-model.number="annonce.poids_estime_kg" type="number" step="0.1" class="input-field" /></div>
            </div>
            <div v-if="annonce.type === 'Vente'" class="form-group"><label class="info-label">Prix</label><input v-model.number="annonce.prix" type="number" step="0.01" class="input-field" /></div>
            <div class="form-group"><label class="info-label">Ville</label><input v-model="annonce.ville" type="text" class="input-field" /></div>
            <div class="specs-grid-edit">
                <div class="form-group"><label class="info-label">Code postal</label><input v-model="annonce.code_postal" type="text" class="input-field" /></div>
                <div class="form-group"><label class="info-label">Adresse</label><input v-model="annonce.adresse" type="text" class="input-field" /></div>
            </div>
        </div>
    </form>
</template>

<script setup>
import { onMounted, ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import { fetchAnnonce, fetchCategories, updateAnnonce } from "../../services/publicApi";

const route = useRoute();
const router = useRouter();
const loading = ref(true);
const submitting = ref(false);
const categories = ref([]);
const annonce = ref({});

const loadPage = async () => {
    loading.value = true;
    try {
        const [categoriesPayload, annoncePayload] = await Promise.all([
            fetchCategories(),
            fetchAnnonce(route.params.id)
        ]);
        categories.value = Array.isArray(categoriesPayload) ? categoriesPayload : [];
        annonce.value = annoncePayload || {};
    } catch (error) {
        console.error("Erreur edition annonce :", error);
    } finally {
        loading.value = false;
    }
};

const handleUpdate = async () => {
    submitting.value = true;
    try {
        await updateAnnonce(annonce.value.id, annonce.value);
        router.push({ name: "see-annonce", params: { id: annonce.value.id } });
    } catch (error) {
        console.error("Erreur mise a jour annonce :", error);
        alert(error.message || "Mise a jour impossible.");
    } finally {
        submitting.value = false;
    }
};

onMounted(loadPage);
</script>

<style scoped>
.split-layout { display: grid; grid-template-columns: 1.5fr 1fr; gap: 2rem; }
.info-card { background: white; padding: 1.5rem; border-radius: 16px; border: 1px solid #eee; margin-bottom: 1.5rem; }
.form-group { margin-bottom: 1.5rem; }
.info-label { display: block; font-size: 0.75rem; font-weight: 700; text-transform: uppercase; color: #aaa; margin-bottom: 8px; letter-spacing: 0.5px; }
.input-field { width: 100%; padding: 0.85rem; border-radius: 10px; border: 1px solid #ddd; background: #fcfcfc; font-size: 1rem; }
.edit-textarea { width: 100%; padding: 1rem; border-radius: 12px; border: 1px solid #eee; background: #f9f9f9; resize: vertical; }
.specs-grid-edit { display: grid; grid-template-columns: 1fr 1fr; gap: 1.5rem; }
.btn-main-action-header { background: #2d6a4f; color: white; padding: 0.6rem 1.5rem; border-radius: 10px; font-weight: bold; cursor: pointer; border: none; }
</style>
