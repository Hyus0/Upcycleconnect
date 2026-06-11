package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"upcycleconnect/api-go/app"
	"upcycleconnect/api-go/db"
)

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Add("school", "esgi")
	res, _ := json.Marshal("en vie.")
	fmt.Fprintf(w, "%s", string(res))
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5174")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	db.Conn = db.NewDB()
	if err := db.EnsureMessagingTables(); err != nil {
		panic(err.Error())
	}
	if err := app.EnsureAdminSchema(); err != nil {
		panic(err.Error())
	}

	http.HandleFunc("GET /", healthCheck)
	http.HandleFunc("GET /health", healthCheck)

	// Admin API consumed by the Vue backoffice.
	http.HandleFunc("GET /api/admin/metrics", app.AdminMetrics)
	http.HandleFunc("GET /api/admin/system/overview", app.AdminSystemOverview)
	http.HandleFunc("GET /api/admin/system/raw", app.AdminRawDump)
	http.HandleFunc("GET /{$}", healthCheck)
	http.HandleFunc("GET /api/admin/users", app.AdminListUsers)
	http.HandleFunc("GET /api/admin/users/{id}", app.AdminGetUser)
	http.HandleFunc("POST /api/admin/users", app.AdminCreateUser)
	http.HandleFunc("PUT /api/admin/users/{id}", app.AdminUpdateUser)
	http.HandleFunc("PATCH /api/admin/users/{id}/status", app.AdminToggleUserStatus)
	http.HandleFunc("DELETE /api/admin/users/{id}", app.AdminDeleteUser)
	http.HandleFunc("GET /api/admin/prestations", app.AdminListPrestations)
	http.HandleFunc("GET /api/admin/prestations/{id}", app.AdminGetPrestation)
	http.HandleFunc("POST /api/admin/prestations", app.AdminCreatePrestation)
	http.HandleFunc("PUT /api/admin/prestations/{id}", app.AdminUpdatePrestation)
	http.HandleFunc("DELETE /api/admin/prestations/{id}", app.AdminDeletePrestation)
	http.HandleFunc("GET /api/admin/categories", app.AdminListCategories)
	http.HandleFunc("GET /api/admin/categories/{id}", app.AdminGetCategory)
	http.HandleFunc("POST /api/admin/categories", app.AdminCreateCategory)
	http.HandleFunc("PUT /api/admin/categories/{id}", app.AdminUpdateCategory)
	http.HandleFunc("DELETE /api/admin/categories/{id}", app.AdminDeleteCategory)
	http.HandleFunc("GET /api/admin/events", app.AdminListEvents)
	http.HandleFunc("GET /api/admin/events/{id}", app.AdminGetEvent)
	http.HandleFunc("POST /api/admin/events", app.AdminCreateEvent)
	http.HandleFunc("PUT /api/admin/events/{id}", app.AdminUpdateEvent)
	http.HandleFunc("DELETE /api/admin/events/{id}", app.AdminDeleteEvent)
	http.HandleFunc("GET /api/admin/moderation/queue", app.AdminModerationQueue)
	http.HandleFunc("PATCH /api/admin/moderation/prestations/{id}/publish", app.AdminPublishPrestation)
	http.HandleFunc("PATCH /api/admin/moderation/prestations/{id}/archive", app.AdminArchivePrestation)
	http.HandleFunc("PATCH /api/admin/moderation/events/{id}/publish", app.AdminPublishEvent)
	http.HandleFunc("PATCH /api/admin/moderation/events/{id}/archive", app.AdminArchiveEvent)
	http.HandleFunc("GET /api/admin/finance/overview", app.AdminFinanceOverview)
	http.HandleFunc("GET /api/admin/subscriptions", app.AdminListSubscriptions)
	http.HandleFunc("GET /api/admin/messages", app.AdminListMessages)
	http.HandleFunc("GET /api/admin/commerce", app.AdminListCommerce)
	http.HandleFunc("GET /api/admin/notifications", app.AdminListNotifications)
	http.HandleFunc("POST /api/admin/notifications", app.AdminCreateNotification)
	http.HandleFunc("PATCH /api/admin/notifications/{id}/status", app.AdminUpdateNotificationStatus)
	http.HandleFunc("DELETE /api/admin/notifications/{id}", app.AdminDeleteNotification)

	// Legacy routes.
	http.HandleFunc("GET /api/admin/user/{id}", app.GetUser)

	//Stat

	http.HandleFunc("GET /stats/platform", app.GetPlatformStatsHandler)

	//Utilisateur
	http.HandleFunc("GET /users", app.GetAllUsers)
	http.HandleFunc("POST /users", app.CreateUser)
	http.HandleFunc("GET /users/{id}", app.GetUser)
	http.HandleFunc("PUT /users/{id}", app.ModifyUser)
	http.HandleFunc("DELETE /users/{id}", app.DeleteUser)
	http.HandleFunc("POST /login", app.UserLogin)
	http.HandleFunc("GET /check-session", app.CheckSession)
	http.HandleFunc("PUT /users/{id}/password", app.ModifyUserPassword)
	http.HandleFunc("GET /users/{id}/stats", app.GetUserStatsHandler)
	http.HandleFunc("GET /user/planning/{id}", app.GetPlanningHandler)
	http.HandleFunc("POST /users/{id}/images", app.UploadUserImages)
	http.Handle("GET /img/", http.StripPrefix("/img/", http.FileServer(http.Dir("./uploads"))))

	// Avis
	http.HandleFunc("GET /users/{id}/avis", app.GetUserAvisHandler)
	http.HandleFunc("POST /users/{id}/avis", app.AddAvisHandler)

	// Follow
	http.HandleFunc("GET /users/{id}/follow/{userId}", app.GetFollowStatusHandler)
	http.HandleFunc("POST /users/{id}/follow/{userId}", app.ToggleFollowHandler)

	//Annonce
	http.HandleFunc("GET /annonces", app.GetAllAnnonces)
	http.HandleFunc("POST /annonces", app.CreateAnnonce)
	http.HandleFunc("GET /annonces/{id}", app.GetAnnonce)
	http.HandleFunc("PUT /annonces/{id}", app.ModifyAnnonce)
	http.HandleFunc("DELETE /annonces/{id}", app.DeleteAnnonce)
	http.HandleFunc("GET /users/{id}/annonces", app.GetUserAnnoncesHandler)
	http.HandleFunc("GET /annonces/{id}/favori/{userId}", app.GetFavoriStatusHandler)
	http.HandleFunc("POST /annonces/{id}/favori/{userId}", app.ToggleFavoriHandler)
	http.HandleFunc("GET /users/{id}/favoris", app.GetMesFavorisHandler)
	http.HandleFunc("GET /users/{id}/achats", app.GetUserAchatsHandler)
	http.HandleFunc("POST /annonces/{id}/image", app.UploadAnnonceImage)

	//Evenement
	http.HandleFunc("GET /evenements", app.GetAllEvenements)
	http.HandleFunc("GET /evenements/{id}", app.GetEvenement)
	http.HandleFunc("POST /evenements", app.CreateEvenement)
	http.HandleFunc("PUT /evenements/{id}", app.ModifyEvenement)
	http.HandleFunc("DELETE /evenements/{id}", app.DeleteEvenement)
	http.HandleFunc("POST /api/evenements/{id}/join", app.JoinEvenement)
	http.HandleFunc("POST /api/evenements/{id}/quit", app.QuitEvenement)
	http.HandleFunc("GET /api/evenements/{id}/inscription-status", app.CheckInscriptionEvenement)
	http.HandleFunc("GET /api/evenements/{id}/participants", app.GetEvenementParticipantsHandler)
	//Categorie
	http.HandleFunc("GET /categories", app.GetAllCategories)
	http.HandleFunc("GET /category/{id}", app.GetCategory)
	http.HandleFunc("POST /category", app.CreateCategory)
	http.HandleFunc("PUT /category/{id}", app.ModifyCategory)
	http.HandleFunc("DELETE /category/{id}", app.DeleteCategory)

	//Formation
	http.HandleFunc("GET /formations", app.GetAllFormations)
	http.HandleFunc("GET /formations/{id}", app.GetFormation)
	http.HandleFunc("POST /formation", app.CreateFormation)
	http.HandleFunc("PUT /formation/{id}", app.ModifyFormation)
	http.HandleFunc("DELETE /formation/{id}", app.DeleteFormation)
	http.HandleFunc("POST /api/formations/{id}/join", app.JoinFormation)
	http.HandleFunc("POST /api/formations/{id}/quit", app.QuitFormation)
	http.HandleFunc("GET /api/formations/{id}/participants", app.GetFormationParticipantsHandler)

	//Logistique
	http.HandleFunc("GET /sites", app.GetAllSites)
	http.HandleFunc("GET /site/{id}", app.GetSiteHandler)
	http.HandleFunc("GET /sites/{id}/conteneurs", app.GetConteneurs)
	http.HandleFunc("POST /annonces/{id}/retirer", app.RetireObjetCasierHandler)
	http.HandleFunc("POST /annonces/{id}/acheter", app.AcheterAnnonceHandler)
	http.HandleFunc("POST /annonces/{id}/reserver", app.ReserverCasierHandler)
	http.HandleFunc("POST /depot", app.DeposerObjetHandler)

	//Projets
	http.HandleFunc("GET /projets", app.GetAllProjets)
	http.HandleFunc("GET /projet/{id}", app.GetProjet)
	http.HandleFunc("POST /projets/{id}/join", app.JoinProjet)
	http.HandleFunc("POST /projets/{id}/quit", app.QuitProjet)
	http.HandleFunc("POST /projets/{id}/like/{userId}", app.ToggleLike)
	http.HandleFunc("GET /projets/{id}/like-status/{userId}", app.CheckLikeStatusHandler)
	http.HandleFunc("GET /users/{id}/projets", app.GetProjetsByUserHandler)
	http.HandleFunc("DELETE /projets/{id}", app.DeleteProjetHandler)
	http.HandleFunc("POST /projets", app.CreateProjetHandler)
	http.HandleFunc("PUT /projets/{id}", app.UpdateProjetHandler)
	http.HandleFunc("POST /projets/upload-image", app.UploadProjetImageHandler)

	//tips
	http.HandleFunc("GET /tips/role/{role}", app.GetTipByRoleHandler)
	http.HandleFunc("GET /tips", app.GetAllTipsHandler)
	http.HandleFunc("GET /tips/{id}", app.GetTipByIDHandler)
	http.HandleFunc("POST /tips", app.CreateTipHandler)
	http.HandleFunc("PUT /tips/{id}", app.UpdateTipHandler)
	http.HandleFunc("DELETE /tips/{id}", app.DeleteTipHandler)

	//commentaires
	http.HandleFunc("GET /commentaires", app.GetAllCommentairesHandler)

	//forums
	http.HandleFunc("GET /forums", app.GetForumsHandler)
	http.HandleFunc("POST /forums/message", app.SendMessageHandler)
	http.HandleFunc("POST /forums/topic", app.CreateTopicHandler)
	http.HandleFunc("DELETE /forums/message/{id}", app.DeleteMessageHandler)
	http.HandleFunc("DELETE /forums/topic/{id}", app.DeleteTopicHandler)
	http.HandleFunc("GET /forums/messages/recent", app.GetRecentMessagesHandler)
	http.HandleFunc("POST /forums/message/{id}/signaler", app.SignalerMessageHandler)
	http.HandleFunc("GET /api/moderation/forums/signales", app.TopMessageSignaleHandler)
	http.HandleFunc("PUT /api/moderation/user/{id}/ban-forum", app.BanUserForumHandler)
	http.HandleFunc("GET /api/moderation/users/banned", app.GetBannedUsersHandler)
	http.HandleFunc("GET /api/moderation/topics", app.GetModerationTopicsHandler)
	//panier
	http.HandleFunc("GET /users/{id}/panier", app.GetPanierHandler)
	http.HandleFunc("POST /users/{id}/panier", app.AddToPanierHandler)
	http.HandleFunc("DELETE /users/{id}/panier/{itemId}", app.RemoveFromPanierHandler)
	http.HandleFunc("POST /users/{id}/checkout", app.CheckoutWithInvoiceHandler)
	http.HandleFunc("GET /users/{id}/factures", app.GetFacturesHandler)
	http.HandleFunc("GET /users/{id}/factures/{factureId}/download", app.DownloadFactureHandler)
	http.HandleFunc("POST /users/{id}/factures/{factureId}/send", app.SendFactureByMailHandler)

	// Messagerie privee
	http.HandleFunc("GET /users/{id}/subscription", app.GetSubscriptionStatusHandler)
	http.HandleFunc("GET /users/{id}/messages", app.GetConversationsHandler)
	http.HandleFunc("POST /users/{id}/messages/start", app.StartConversationHandler)
	http.HandleFunc("GET /users/{id}/messages/{conversationId}", app.GetConversationMessagesHandler)
	http.HandleFunc("POST /users/{id}/messages/{conversationId}", app.SendDMMessageHandler)
	http.HandleFunc("GET /users/{id}/messages/{conversationId}/state", app.GetConversationStateHandler)
	http.HandleFunc("POST /users/{id}/messages/{conversationId}/offers", app.CreateDMOfferHandler)
	http.HandleFunc("PATCH /users/{id}/messages/offers/{offerId}", app.RespondDMOfferHandler)
	http.HandleFunc("POST /users/{id}/messages/sales/{saleId}/reception", app.ConfirmDMSaleReceptionHandler)
	http.HandleFunc("POST /users/{id}/messages/sales/{saleId}/review", app.ReviewDMSaleHandler)

	//notification
	http.HandleFunc("GET /notifications", app.GetAllNotificationsHandler)
	http.HandleFunc("GET /notifications/{id}", app.GetNotificationHandler)
	http.HandleFunc("GET /users/{id}/notifications", app.GetNotificationsHandler)
	http.HandleFunc("POST /notifications/{id}/read", app.MarquerNotificationLueHandler)

	//traduction
	http.HandleFunc("GET /langues", app.GetLanguesHandler)
	http.HandleFunc("GET /traductions/{code}", app.GetTraductionsHandler)
	http.HandleFunc("PUT /users/{id}/langue", app.UpdateLangueHandler)

	fmt.Println("Listening at http://localhost:8081")
	http.ListenAndServe(":8081", enableCORS(http.DefaultServeMux))
}
