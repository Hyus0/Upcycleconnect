<template>
    <div class="info-view">
        <header class="content-header">
            <div class="header-left">
                <p class="sidebar__category2">ACCUEIL &gt; PARAMETRES &gt; INFORMATIONS</p>
                <h1 class="hero-title1">MON PROFIL</h1>
                <p class="classic-text">Les informations de compte modifiees ici viennent et repartent vers le back.</p>
                <div v-if="errors.length > 0" class="error-box">
                    <ul><li v-for="(err, index) in errors" :key="index">{{ err }}</li></ul>
                </div>
                <div v-if="successMsg" class="success-box">{{ successMsg }}</div>
            </div>
            <button class="btn-main-action" @click="updateProfile">Mettre a jour le profil</button>
        </header>

        <div class="info-layout">
            <div class="info-main-col">
                <section class="info-section">
                    <h2 class="section-title">Identite personnelle</h2>
                    <div class="input-grid">
                        <div class="input-field"><label>Prenom</label><input v-model="form.prenom" type="text" /></div>
                        <div class="input-field"><label>Nom</label><input v-model="form.nom" type="text" /></div>
                        <div class="input-field"><label>Date de naissance</label><input v-model="form.date_naissance" type="date" /></div>
                        <div class="input-field"><label>Email</label><input v-model="form.mail" type="email" disabled class="disabled-input" /></div>
                    </div>
                </section>

                <section class="info-section">
                    <h2 class="section-title">Adresse et localisation</h2>
                    <div class="input-field"><label>Adresse complete</label><input v-model="form.adresse" type="text" /></div>
                    <div class="input-grid">
                        <div class="input-field"><label>Ville</label><input v-model="form.ville" type="text" /></div>
                        <div class="input-field"><label>Code postal</label><input v-model="form.code_postal" type="text" /></div>
                    </div>
                </section>
            </div>

            <aside class="info-side-col">
                <section class="info-section status-card">
                    <h2 class="section-title">Etat du compte</h2>
                    <div class="status-item"><span class="status-label">Role</span><span class="badge-role">{{ form.role || "NULL" }}</span></div>
                    <div class="status-item"><span class="status-label">Statut</span><span class="status-active">{{ form.statut || "NULL" }}</span></div>
                    <div class="status-item"><span class="status-label">Langue ID</span><span>{{ form.id_langue ?? "NULL" }}</span></div>
                    <div class="divider"></div>
                    <div class="registration-info">
                        <p>Membre depuis le :</p>
                        <strong>{{ formatDate(form.date_inscription) }}</strong>
                    </div>
                </section>

                <section class="info-section danger-zone">
                    <h2 class="section-title">Securite</h2>
                    <button class="btn-outline" @click="router.push('/profil/password')">Changer le mot de passe</button>
                </section>
            </aside>
        </div>
    </div>
</template>

<script setup>
import { onMounted, ref } from "vue";
import { useRouter } from "vue-router";
import { fetchUser, updateUser } from "../../services/publicApi";

const router = useRouter();
const form = ref({
    prenom: "",
    nom: "",
    mail: "",
    adresse: "",
    ville: "",
    code_postal: "",
    date_naissance: "",
    date_inscription: "",
    role: "",
    statut: "",
    id_langue: null
});
const errors = ref([]);
const successMsg = ref("");

const formatDate = (dateString) => {
    if (!dateString) return "NULL";
    return new Date(dateString).toLocaleDateString("fr-FR", {
        day: "numeric",
        month: "long",
        year: "numeric"
    });
};

const loadProfile = async () => {
    const id = localStorage.getItem("userId");
    if (!id) return;
    try {
        form.value = await fetchUser(id);
    } catch (error) {
        console.error("Erreur profil :", error);
    }
};

const updateProfile = async () => {
    const userId = localStorage.getItem("userId");
    errors.value = [];
    successMsg.value = "";
    if (!userId) return;

    try {
        await updateUser(userId, form.value);
        successMsg.value = "Profil mis a jour avec succes.";
        localStorage.setItem("userPrenom", form.value.prenom || "");
        localStorage.setItem("userNom", form.value.nom || "");
        window.dispatchEvent(new Event("auth-change"));
    } catch (error) {
        errors.value = [error.message || "Une erreur est survenue."];
    }
};

onMounted(loadProfile);
</script>

<style scoped>
.info-view { padding: 10px; }
.info-layout { display: grid; grid-template-columns: 1fr 320px; gap: 25px; margin-top: 20px; }
.info-section { background: white; border-radius: 16px; padding: 24px; margin-bottom: 25px; border: 1px solid #e8ebe9; }
.section-title { font-size: 1.1rem; font-weight: 700; margin-bottom: 20px; color: #1a1f1c; }
.input-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 20px; }
.input-field { display: flex; flex-direction: column; gap: 8px; margin-bottom: 15px; }
label { font-size: 0.85rem; font-weight: 600; color: #5a6660; }
input { background-color: #fcfdfc; border: 1px solid #dcdfdc; border-radius: 10px; padding: 12px 15px; font-size: 0.95rem; }
.disabled-input { background-color: #f1f3f2; color: #888; cursor: not-allowed; }
.status-item { display: flex; justify-content: space-between; margin-bottom: 15px; font-size: 0.9rem; }
.badge-role { background-color: #eaf4ed; color: #2d7a4f; padding: 4px 10px; border-radius: 8px; font-weight: bold; font-size: 0.8rem; }
.status-active { color: #1e7e34; font-weight: bold; }
.divider { height: 1px; background: #e8ebe9; margin: 20px 0; }
.btn-outline { width: 100%; background: white; border: 1px solid #dcdfdc; padding: 10px; border-radius: 10px; cursor: pointer; font-weight: 600; }
.error-box { background-color: #fee2e2; border: 1px solid #ef4444; color: #b91c1c; padding: 10px; border-radius: 8px; margin-bottom: 15px; }
.success-box { background-color: #e2fee3; border: 1px solid #44ef44; color: #158f3c; padding: 10px; border-radius: 8px; margin-bottom: 15px; }
</style>
