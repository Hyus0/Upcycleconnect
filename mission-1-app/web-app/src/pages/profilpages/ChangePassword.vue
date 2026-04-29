<template>
    <div class="info-view">
        <header class="content-header">
            <div class="header-left">
                <p class="sidebar__category2">ACCUEIL &gt; PARAMETRES &gt; MOT DE PASSE</p>
                <h1 class="hero-title1">MODIFIER MON MOT DE PASSE</h1>
                <p class="classic-text">La mise a jour passe par l'API compte et est enregistree en base.</p>
                <div v-if="errors.length" class="error-box"><ul><li v-for="(err, index) in errors" :key="index">{{ err }}</li></ul></div>
                <div v-if="successMsg" class="success-box">{{ successMsg }}</div>
            </div>
            <button class="btn-main-action" @click="handlePasswordUpdate">Enregistrer le mot de passe</button>
        </header>

        <div class="info-layout">
            <section class="info-section">
                <div class="input-field"><label>Mot de passe actuel</label><input v-model="pwdForm.oldPassword" type="password" /></div>
                <div class="divider"></div>
                <div class="input-grid">
                    <div class="input-field"><label>Nouveau mot de passe</label><input v-model="pwdForm.newPassword" type="password" /></div>
                    <div class="input-field"><label>Confirmation</label><input v-model="pwdForm.confirmPassword" type="password" /></div>
                </div>
            </section>

            <aside class="info-section status-card">
                <h2 class="section-title">Derniere modification</h2>
                <p><strong>{{ formatDate(lastChangedDate) }}</strong></p>
            </aside>
        </div>
    </div>
</template>

<script setup>
import { onMounted, ref } from "vue";
import { fetchUser, updatePassword } from "../../services/publicApi";

const errors = ref([]);
const successMsg = ref("");
const lastChangedDate = ref("");
const pwdForm = ref({
    oldPassword: "",
    newPassword: "",
    confirmPassword: ""
});

const formatDate = (value) => {
    if (!value || value === "0000-00-00 00:00:00") return "Jamais";
    const date = new Date(value);
    return Number.isNaN(date.getTime())
        ? "Format invalide"
        : date.toLocaleDateString("fr-FR", { day: "numeric", month: "long", year: "numeric" });
};

const loadProfile = async () => {
    const id = localStorage.getItem("userId");
    if (!id) return;
    try {
        const payload = await fetchUser(id);
        lastChangedDate.value = payload?.date_update_password || "";
    } catch (error) {
        console.error("Erreur chargement profil :", error);
    }
};

const handlePasswordUpdate = async () => {
    errors.value = [];
    successMsg.value = "";
    if (pwdForm.value.newPassword !== pwdForm.value.confirmPassword) {
        errors.value.push("Le nouveau mot de passe et la confirmation ne correspondent pas.");
        return;
    }
    if (pwdForm.value.newPassword.length < 8) {
        errors.value.push("Le nouveau mot de passe doit contenir au moins 8 caracteres.");
        return;
    }

    const userId = localStorage.getItem("userId");
    if (!userId) return;

    try {
        await updatePassword(userId, {
            old_password: pwdForm.value.oldPassword,
            password: pwdForm.value.newPassword
        });
        successMsg.value = "Votre mot de passe a ete mis a jour.";
        lastChangedDate.value = new Date().toISOString();
        pwdForm.value = { oldPassword: "", newPassword: "", confirmPassword: "" };
    } catch (error) {
        errors.value = [error.message || "Le serveur est injoignable."];
    }
};

onMounted(loadProfile);
</script>

<style scoped>
.info-view { padding: 10px; }
.info-layout { display: grid; grid-template-columns: 1fr 320px; gap: 25px; margin-top: 20px; }
.info-section { background: white; border-radius: 16px; padding: 24px; margin-bottom: 25px; border: 1px solid #e8ebe9; }
.input-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 20px; }
.input-field { display: flex; flex-direction: column; gap: 8px; margin-bottom: 15px; }
input { background-color: #fcfdfc; border: 1px solid #dcdfdc; border-radius: 10px; padding: 12px 15px; font-size: 0.95rem; }
.divider { height: 1px; background: #e8ebe9; margin: 20px 0; }
.error-box { background-color: #fee2e2; border: 1px solid #ef4444; color: #b91c1c; padding: 12px; border-radius: 10px; margin-bottom: 15px; }
.success-box { background-color: #f0fdf4; border: 1px solid #22c55e; color: #166534; padding: 12px; border-radius: 10px; margin-bottom: 15px; }
</style>
