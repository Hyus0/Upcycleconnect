<template>
    <div class="info-view">
        <header class="content-header">
            <div class="header-left">
                <p class="sidebar__category2">
                    ACCUEIL > PARAMÈTRES > INFORMATIONS > CHANGER MON MOT DE PASSE
                </p>
                <h1 class="hero-title1">MODIFIER MON MOT DE PASSE</h1>
                <p class="classic-text">
                    Pour des raisons de sécurité, votre mot de passe doit être complexe et ne pas être réutilisé pour d'autres comptes.
                </p>
                
                <div v-if="errors.length > 0" class="error-box">
                    <ul style="margin: 0; padding-left: 20px;">
                        <li v-for="(err, index) in errors" :key="index">{{ err }}</li>
                    </ul>
                </div>
                <div v-if="successMsg" class="success-box">
                    {{ successMsg }}
                </div>
            </div>
            <button class="btn-main-action" @click="handlePasswordUpdate">
                Enregistrer le mot de passe
            </button>
        </header>

        <div class="info-layout">
            <div class="info-main-col">
                <section class="info-section">
                    <h2 class="section-title">Sécurité du compte</h2>
                    <div class="password-form">
                        <div class="input-field">
                            <label>Mot de passe actuel</label>
                            <input 
                                type="password" 
                                v-model="pwdForm.oldPassword" 
                                placeholder="••••••••"
                            />
                        </div>
                        
                        <div class="divider"></div>

                        <div class="input-grid">
                            <div class="input-field">
                                <label>Nouveau mot de passe</label>
                                <input 
                                    type="password" 
                                    v-model="pwdForm.newPassword" 
                                    placeholder="Minimum 8 caractères"
                                />
                            </div>
                            <div class="input-field">
                                <label>Confirmer le nouveau mot de passe</label>
                                <input 
                                    type="password" 
                                    v-model="pwdForm.confirmPassword" 
                                    placeholder="••••••••"
                                />
                            </div>
                        </div>
                    </div>
                </section>
                
                <section class="info-section tips-card">
                    <h3 class="section-title" style="font-size: 0.9rem;">💡 Conseils pour un mot de passe robuste</h3>
                    <ul class="tips-list">
                        <li>Utilisez au moins 8 caractères.</li>
                        <li>Mélangez majuscules, minuscules et chiffres.</li>
                        <li>Ajoutez un caractère spécial (ex: @, #, $, !).</li>
                    </ul>
                </section>
            </div>

            <aside class="info-side-col">
                <section class="info-section status-card">
                    <h2 class="section-title">Dernière modification</h2>
                    <p style="font-size: 0.85rem; color: #5a6660;">
                        Votre mot de passe a été modifié pour la dernière fois le :<br>
                        <strong>{{ formatDate(lastChangedDate) }}</strong>
                    </p>
                </section>
            </aside>
        </div>
    </div>
</template>

<script setup>
import { ref, onMounted } from "vue";

const errors = ref([]);
const successMsg = ref("");

const lastChangedDate = ref("");

const pwdForm = ref({
    oldPassword: "",
    newPassword: "",
    confirmPassword: "",
});

const formatDate = (dateString) => {
    if (!dateString || dateString === "" || dateString === "0000-00-00 00:00:00") {
        return "Jamais";
    }
    
    const date = new Date(dateString);
    
    if (isNaN(date.getTime())) return "Format invalide";

    return date.toLocaleDateString("fr-FR", {
        day: "numeric",
        month: "long",
        year: "numeric",
    });
};

onMounted(async () => {
    const id = localStorage.getItem("userId");
    const token = localStorage.getItem("userToken");
    if (!id || !token) return;

    try {
        const response = await fetch(`http://localhost:8081/users/${id}`, {
            headers: { Authorization: token }
        });
        if (response.ok) {
            const data = await response.json();
            lastChangedDate.value = data.date_update_password;
        }
    } catch (err) {
        console.error("Erreur chargement profil :", err);
    }
});

const handlePasswordUpdate = async () => {
    errors.value = [];
    successMsg.value = "";

    if (pwdForm.value.newPassword !== pwdForm.value.confirmPassword) {
        errors.value.push("Le nouveau mot de passe et la confirmation ne correspondent pas.");
        return;
    }

    if (pwdForm.value.newPassword.length < 8) {
        errors.value.push("Le nouveau mot de passe doit contenir au moins 8 caractères.");
        return;
    }

    const userId = localStorage.getItem("userId");
    const token = localStorage.getItem("userToken");

    try {
        const response = await fetch(`http://localhost:8081/users/${userId}/password`, {
            method: "PUT",
            headers: {
                "Authorization": token,
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                old_password: pwdForm.value.oldPassword,
                password: pwdForm.value.newPassword
            }),
        });

        if (response.ok) {
            successMsg.value = "Votre mot de passe a été mis à jour ! 🎉";
            lastChangedDate.value = new Date().toISOString(); 
            pwdForm.value = { oldPassword: "", newPassword: "", confirmPassword: "" };
        } else {
            const data = await response.json();
            errors.value = Array.isArray(data) ? data : [data.message || "Erreur lors de la mise à jour"];
        }
    } catch (error) {
        errors.value = ["Le serveur est injoignable."];
    }
};
</script>

<style scoped>
.info-view {
    padding: 10px;
}

.info-layout {
    display: grid;
    grid-template-columns: 1fr 320px;
    gap: 25px;
    margin-top: 20px;
}

.info-section {
    background: white;
    border-radius: 16px;
    padding: 24px;
    margin-bottom: 25px;
    border: 1px solid #e8ebe9;
}

.section-title {
    font-size: 1.1rem;
    font-family: "Syne", sans-serif;
    font-weight: 700;
    margin-bottom: 20px;
    color: #1a1f1c;
}

.password-form {
    display: flex;
    flex-direction: column;
    gap: 10px;
}

.input-grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 20px;
}

.input-field {
    display: flex;
    flex-direction: column;
    gap: 8px;
    margin-bottom: 15px;
}

label {
    font-size: 0.85rem;
    font-weight: 600;
    color: #5a6660;
}

input {
    background-color: #fcfdfc;
    border: 1px solid #dcdfdc;
    border-radius: 10px;
    padding: 12px 15px;
    font-size: 0.95rem;
    outline: none;
    transition: all 0.2s;
}

input:focus {
    border-color: #2d7a4f;
    background-color: #fff;
    box-shadow: 0 0 0 3px rgba(45, 122, 79, 0.1);
}

.divider {
    height: 1px;
    background: #e8ebe9;
    margin: 20px 0;
}

.error-box {
    background-color: #fee2e2;
    border: 1px solid #ef4444;
    color: #b91c1c;
    padding: 12px;
    border-radius: 10px;
    margin-bottom: 15px;
    font-size: 0.9rem;
}

.success-box {
    background-color: #f0fdf4;
    border: 1px solid #22c55e;
    color: #166534;
    padding: 12px;
    border-radius: 10px;
    margin-bottom: 15px;
    font-size: 0.9rem;
}

.tips-card {
    background-color: #f8faf9;
    border: 1px dashed #ced6d1;
}

.tips-list {
    font-size: 0.85rem;
    color: #5a6660;
    line-height: 1.6;
    padding-left: 20px;
    margin: 0;
}

.status-card {
    height: fit-content;
}

/* --- RESPONSIVE --- */
@media (max-width: 1024px) {
    .info-layout {
        grid-template-columns: 1fr;
    }
    .info-side-col {
        order: -1;
    }
    .input-grid {
        grid-template-columns: 1fr;
    }
}
</style>