package i18n

// String keys for internationalization.
const (
	// Button strings.
	ButtonStartConversation = "button.start_conversation"
	ButtonModelSelect       = "button.model_select"
	ButtonBackToMenu        = "button.back_to_menu"
	ButtonNewConversation   = "button.new_conversation"
	ButtonSettings          = "button.settings"
	ButtonLanguage          = "button.language"
	ButtonPrevPage          = "button.prev_page"
	ButtonNextPage          = "button.next_page"

	// Menu messages.
	MenuWelcome    = "menu.welcome"
	MenuBackToMain = "menu.back_to_main"

	// Conversation messages.
	ConversationStarted       = "conversation.started"
	ConversationResumed       = "conversation.resumed"
	ConversationModePrompt    = "conversation.mode_prompt"
	ConversationNameGenerated = "conversation.name_generated"

	// Settings messages.
	SettingsTitle         = "settings.title"
	LanguageSelectTitle   = "settings.language_select_title"
	LanguageUpdateSuccess = "settings.language_update_success"

	// Conversation list messages.
	ConversationListSelect   = "conversation_list.select"
	ConversationListEmpty    = "conversation_list.empty"
	ConversationListPageInfo = "conversation_list.page_info"

	// Model selection messages.
	ModelSelectTitle   = "model.select_title"
	ModelUpdateSuccess = "model.update_success"

	// Queue messages.
	QueueMessageQueued = "queue.message_queued"

	// Language names (for language selection).
	LangEnglish    = "lang.english"
	LangSpanish    = "lang.spanish"
	LangFrench     = "lang.french"
	LangGerman     = "lang.german"
	LangItalian    = "lang.italian"
	LangRussian    = "lang.russian"
	LangChinese    = "lang.chinese"
	LangJapanese   = "lang.japanese"
	LangKorean     = "lang.korean"
	LangPortuguese = "lang.portuguese"
	LangArmenian   = "lang.armenian"
	LangUkrainian  = "lang.ukrainian"
	LangKazakh     = "lang.kazakh"
	LangKyrgyz     = "lang.kyrgyz"
	LangArabic     = "lang.arabic"
	LangHindi      = "lang.hindi"
)

// Languages maps language codes to their display names with flags.
var Languages = map[string]string{
	LangEnglish:    "🇺🇸 English",
	LangSpanish:    "🇪🇸 Español",
	LangFrench:     "🇫🇷 Français",
	LangGerman:     "🇩🇪 Deutsch",
	LangItalian:    "🇮🇹 Italiano",
	LangRussian:    "🇷🇺 Русский",
	LangChinese:    "🇨🇳 中文",
	LangJapanese:   "🇯🇵 日本語",
	LangKorean:     "🇰🇷 한국어",
	LangPortuguese: "🇵🇹 Português",
	LangArmenian:   "🇦🇲 Հայերեն",
	LangUkrainian:  "🇺🇦 Українська",
	LangKazakh:     "🇰🇿 Қазақша",
	LangKyrgyz:     "🇰🇬 Кыргызча",
	LangArabic:     "🇸🇦 العربية",
	LangHindi:      "🇮🇳 हिन्दी",
}

// LanguageCodes maps language names to their 2-letter codes.
var LanguageCodes = map[string]string{
	Languages[LangEnglish]:    "en",
	Languages[LangSpanish]:    "es",
	Languages[LangFrench]:     "fr",
	Languages[LangGerman]:     "de",
	Languages[LangItalian]:    "it",
	Languages[LangRussian]:    "ru",
	Languages[LangChinese]:    "zh",
	Languages[LangJapanese]:   "ja",
	Languages[LangKorean]:     "ko",
	Languages[LangPortuguese]: "pt",
	Languages[LangArmenian]:   "hy",
	Languages[LangUkrainian]:  "uk",
	Languages[LangKazakh]:     "kk",
	Languages[LangKyrgyz]:     "ky",
	Languages[LangArabic]:     "ar",
	Languages[LangHindi]:      "hi",
}

// Strings holds all localized strings.
var Strings = map[string]map[string]string{
	"en": {
		// Buttons
		ButtonStartConversation: "💬️ Start Conversation",
		ButtonModelSelect:       "🤖 Select Model",
		ButtonBackToMenu:        "🔙 Back to Menu",
		ButtonNewConversation:   "➕ New Conversation",
		ButtonSettings:          "⚙️ Settings",
		ButtonLanguage:          "🌐 Language",
		ButtonPrevPage:          "⬅️",
		ButtonNextPage:          "➡️",

		// Menu
		MenuWelcome:    "Welcome! Choose an option:",
		MenuBackToMain: "🏠 Back to main menu. Choose an option:",

		// Conversation
		ConversationStarted:       "🗣️ Conversation started! Send me a message and I'll respond. You can always go back to the menu.",
		ConversationResumed:       "🗣️ Conversation resumed! Send me a message and I'll respond. You can always go back to the menu.",
		ConversationModePrompt:    "You're in conversation mode. Send a message to chat, or go back to menu:",
		ConversationNameGenerated: "conversation name generated successfully",

		// Settings
		SettingsTitle:         "⚙️ Settings. Choose an option:",
		LanguageSelectTitle:   "🌐 Select your language:",
		LanguageUpdateSuccess: "✅ Language updated successfully!",

		// Conversation List
		ConversationListSelect:   "💬 Select a conversation:",
		ConversationListEmpty:    "💬 No previous conversations. Start a new one:",
		ConversationListPageInfo: "💬 Select a conversation (page %d/%d):",

		// Model Selection
		ModelSelectTitle:   "🤖 Select AI model:",
		ModelUpdateSuccess: "✅ Model updated successfully!",

		// Queue
		QueueMessageQueued: "⏳ Your message has been queued (position: %d). I'll process it after finishing the current response.",

		// Languages
		LangEnglish:    "🇺🇸 English",
		LangSpanish:    "🇪🇸 Español",
		LangFrench:     "🇫🇷 Français",
		LangGerman:     "🇩🇪 Deutsch",
		LangItalian:    "🇮🇹 Italiano",
		LangRussian:    "🇷🇺 Русский",
		LangChinese:    "🇨🇳 中文",
		LangJapanese:   "🇯🇵 日本語",
		LangKorean:     "🇰🇷 한국어",
		LangPortuguese: "🇵🇹 Português",
		LangArmenian:   "🇦🇲 Հայերեն",
		LangUkrainian:  "🇺🇦 Українська",
		LangKazakh:     "🇰🇿 Қазақша",
		LangKyrgyz:     "🇰🇬 Кыргызча",
		LangArabic:     "🇸🇦 العربية",
		LangHindi:      "🇮🇳 हिन्दी",
	},
	"es": {
		// Buttons
		ButtonStartConversation: "💬️ Iniciar Conversación",
		ButtonModelSelect:       "🤖 Seleccionar Modelo",
		ButtonBackToMenu:        "🔙 Volver al Menú",
		ButtonNewConversation:   "➕ Nueva Conversación",
		ButtonSettings:          "⚙️ Configuración",
		ButtonLanguage:          "🌐 Idioma",
		ButtonPrevPage:          "⬅️",
		ButtonNextPage:          "➡️",

		// Menu
		MenuWelcome:    "¡Bienvenido! Elige una opción:",
		MenuBackToMain: "🏠 Volver al menú principal. Elige una opción:",

		// Conversation
		ConversationStarted:       "🗣️ ¡Conversación iniciada! Envíame un mensaje y te responderé. Siempre puedes volver al menú.",
		ConversationResumed:       "🗣️ ¡Conversación reanudada! Envíame un mensaje y te responderé. Siempre puedes volver al menú.",
		ConversationModePrompt:    "Estás en modo conversación. Envía un mensaje para chatear, o vuelve al menú:",
		ConversationNameGenerated: "nombre de conversación generado exitosamente",

		// Settings
		SettingsTitle:         "⚙️ Configuración. Elige una opción:",
		LanguageSelectTitle:   "🌐 Selecciona tu idioma:",
		LanguageUpdateSuccess: "✅ ¡Idioma actualizado exitosamente!",

		// Conversation List
		ConversationListSelect:   "💬 Selecciona una conversación:",
		ConversationListEmpty:    "💬 No hay conversaciones anteriores. Inicia una nueva:",
		ConversationListPageInfo: "💬 Selecciona una conversación (página %d/%d):",

		// Model Selection
		ModelSelectTitle:   "🤖 Seleccionar modelo de IA:",
		ModelUpdateSuccess: "✅ ¡Modelo actualizado exitosamente!",

		// Queue
		QueueMessageQueued: "⏳ Tu mensaje ha sido puesto en cola (posición: %d). Lo procesaré después de terminar la respuesta actual.",

		// Languages
		LangEnglish:    "🇺🇸 English",
		LangSpanish:    "🇪🇸 Español",
		LangFrench:     "🇫🇷 Français",
		LangGerman:     "🇩🇪 Deutsch",
		LangItalian:    "🇮🇹 Italiano",
		LangRussian:    "🇷🇺 Русский",
		LangChinese:    "🇨🇳 中文",
		LangJapanese:   "🇯🇵 日本語",
		LangKorean:     "🇰🇷 한국어",
		LangPortuguese: "🇵🇹 Português",
		LangArmenian:   "🇦🇲 Հայերեն",
		LangUkrainian:  "🇺🇦 Українська",
		LangKazakh:     "🇰🇿 Қазақша",
		LangKyrgyz:     "🇰🇬 Кыргызча",
		LangArabic:     "🇸🇦 العربية",
		LangHindi:      "🇮🇳 हिन्दी",
	},
	"ru": {
		// Buttons
		ButtonStartConversation: "💬️ Начать Беседу",
		ButtonModelSelect:       "🤖 Выбрать Модель",
		ButtonBackToMenu:        "🔙 Назад в Меню",
		ButtonNewConversation:   "➕ Новая Беседа",
		ButtonSettings:          "⚙️ Настройки",
		ButtonLanguage:          "🌐 Язык",
		ButtonPrevPage:          "⬅️",
		ButtonNextPage:          "➡️",

		// Menu
		MenuWelcome:    "Добро пожаловать! Выберите опцию:",
		MenuBackToMain: "🏠 Назад в главное меню. Выберите опцию:",

		// Conversation
		ConversationStarted:       "🗣️ Беседа начата! Отправьте мне сообщение, и я отвечу. Вы всегда можете вернуться в меню.",
		ConversationResumed:       "🗣️ Беседа возобновлена! Отправьте мне сообщение, и я отвечу. Вы всегда можете вернуться в меню.",
		ConversationModePrompt:    "Вы в режиме беседы. Отправьте сообщение для чата или вернитесь в меню:",
		ConversationNameGenerated: "название беседы успешно создано",

		// Settings
		SettingsTitle:         "⚙️ Настройки. Выберите опцию:",
		LanguageSelectTitle:   "🌐 Выберите ваш язык:",
		LanguageUpdateSuccess: "✅ Язык успешно обновлен!",

		// Conversation List
		ConversationListSelect:   "💬 Выберите беседу:",
		ConversationListEmpty:    "💬 Нет предыдущих бесед. Начните новую:",
		ConversationListPageInfo: "💬 Выберите беседу (страница %d/%d):",

		// Model Selection
		ModelSelectTitle:   "🤖 Выбрать ИИ модель:",
		ModelUpdateSuccess: "✅ Модель успешно обновлена!",

		// Queue
		QueueMessageQueued: "⏳ Ваше сообщение поставлено в очередь (позиция: %d). Я обработаю его после завершения текущего ответа.",

		// Languages
		LangEnglish:    "🇺🇸 English",
		LangSpanish:    "🇪🇸 Español",
		LangFrench:     "🇫🇷 Français",
		LangGerman:     "🇩🇪 Deutsch",
		LangItalian:    "🇮🇹 Italiano",
		LangRussian:    "🇷🇺 Русский",
		LangChinese:    "🇨🇳 中文",
		LangJapanese:   "🇯🇵 日本語",
		LangKorean:     "🇰🇷 한국어",
		LangPortuguese: "🇵🇹 Português",
		LangArmenian:   "🇦🇲 Հայերեն",
		LangUkrainian:  "🇺🇦 Українська",
		LangKazakh:     "🇰🇿 Қазақша",
		LangKyrgyz:     "🇰🇬 Кыргызча",
		LangArabic:     "🇸🇦 العربية",
		LangHindi:      "🇮🇳 हिन्दी",
	},
	"fr": {
		// Buttons
		ButtonStartConversation: "💬️ Commencer une Conversation",
		ButtonModelSelect:       "🤖 Sélectionner un Modèle",
		ButtonBackToMenu:        "🔙 Retour au Menu",
		ButtonNewConversation:   "➕ Nouvelle Conversation",
		ButtonSettings:          "⚙️ Paramètres",
		ButtonLanguage:          "🌐 Langue",
		ButtonPrevPage:          "⬅️",
		ButtonNextPage:          "➡️",

		// Menu
		MenuWelcome:    "Bienvenue ! Choisissez une option :",
		MenuBackToMain: "🏠 Retour au menu principal. Choisissez une option :",

		// Conversation
		ConversationStarted:       "🗣️ Conversation commencée ! Envoyez-moi un message et je répondrai. Vous pouvez toujours retourner au menu.",
		ConversationResumed:       "🗣️ Conversation reprise ! Envoyez-moi un message et je répondrai. Vous pouvez toujours retourner au menu.",
		ConversationModePrompt:    "Vous êtes en mode conversation. Envoyez un message pour discuter, ou retournez au menu :",
		ConversationNameGenerated: "nom de conversation généré avec succès",

		// Settings
		SettingsTitle:         "⚙️ Paramètres. Choisissez une option :",
		LanguageSelectTitle:   "🌐 Sélectionnez votre langue :",
		LanguageUpdateSuccess: "✅ Langue mise à jour avec succès !",

		// Conversation List
		ConversationListSelect:   "💬 Sélectionnez une conversation :",
		ConversationListEmpty:    "💬 Aucune conversation précédente. Commencez-en une nouvelle :",
		ConversationListPageInfo: "💬 Sélectionnez une conversation (page %d/%d) :",

		// Model Selection
		ModelSelectTitle:   "🤖 Sélectionner un modèle IA :",
		ModelUpdateSuccess: "✅ Modèle mis à jour avec succès !",

		// Queue
		QueueMessageQueued: "⏳ Votre message a été mis en file d'attente (position : %d). Je le traiterai après avoir terminé la réponse actuelle.",

		// Languages
		LangEnglish:    "🇺🇸 English",
		LangSpanish:    "🇪🇸 Español",
		LangFrench:     "🇫🇷 Français",
		LangGerman:     "🇩🇪 Deutsch",
		LangItalian:    "🇮🇹 Italiano",
		LangRussian:    "🇷🇺 Русский",
		LangChinese:    "🇨🇳 中文",
		LangJapanese:   "🇯🇵 日本語",
		LangKorean:     "🇰🇷 한국어",
		LangPortuguese: "🇵🇹 Português",
		LangArmenian:   "🇦🇲 Հայերեն",
		LangUkrainian:  "🇺🇦 Українська",
		LangKazakh:     "🇰🇿 Қазақша",
		LangKyrgyz:     "🇰🇬 Кыргызча",
		LangArabic:     "🇸🇦 العربية",
		LangHindi:      "🇮🇳 हिन्दी",
	},
	"de": {
		// Buttons
		ButtonStartConversation: "💬️ Gespräch Starten",
		ButtonModelSelect:       "🤖 Modell Auswählen",
		ButtonBackToMenu:        "🔙 Zurück zum Menü",
		ButtonNewConversation:   "➕ Neues Gespräch",
		ButtonSettings:          "⚙️ Einstellungen",
		ButtonLanguage:          "🌐 Sprache",
		ButtonPrevPage:          "⬅️",
		ButtonNextPage:          "➡️",

		// Menu
		MenuWelcome:    "Willkommen! Wählen Sie eine Option:",
		MenuBackToMain: "🏠 Zurück zum Hauptmenü. Wählen Sie eine Option:",

		// Conversation
		ConversationStarted:       "🗣️ Gespräch gestartet! Senden Sie mir eine Nachricht und ich werde antworten. Sie können jederzeit zum Menü zurückkehren.",
		ConversationResumed:       "🗣️ Gespräch fortgesetzt! Senden Sie mir eine Nachricht und ich werde antworten. Sie können jederzeit zum Menü zurückkehren.",
		ConversationModePrompt:    "Sie befinden sich im Gesprächsmodus. Senden Sie eine Nachricht zum Chatten oder kehren Sie zum Menü zurück:",
		ConversationNameGenerated: "Gesprächsname erfolgreich generiert",

		// Settings
		SettingsTitle:         "⚙️ Einstellungen. Wählen Sie eine Option:",
		LanguageSelectTitle:   "🌐 Wählen Sie Ihre Sprache:",
		LanguageUpdateSuccess: "✅ Sprache erfolgreich aktualisiert!",

		// Conversation List
		ConversationListSelect:   "💬 Wählen Sie ein Gespräch:",
		ConversationListEmpty:    "💬 Keine vorherigen Gespräche. Starten Sie ein neues:",
		ConversationListPageInfo: "💬 Wählen Sie ein Gespräch (Seite %d/%d):",

		// Model Selection
		ModelSelectTitle:   "🤖 KI-Modell auswählen:",
		ModelUpdateSuccess: "✅ Modell erfolgreich aktualisiert!",

		// Queue
		QueueMessageQueued: "⏳ Ihre Nachricht wurde in die Warteschlange eingereiht (Position: %d). Ich werde sie nach Beendigung der aktuellen Antwort bearbeiten.",

		// Languages
		LangEnglish:    "🇺🇸 English",
		LangSpanish:    "🇪🇸 Español",
		LangFrench:     "🇫🇷 Français",
		LangGerman:     "🇩🇪 Deutsch",
		LangItalian:    "🇮🇹 Italiano",
		LangRussian:    "🇷🇺 Русский",
		LangChinese:    "🇨🇳 中文",
		LangJapanese:   "🇯🇵 日本語",
		LangKorean:     "🇰🇷 한국어",
		LangPortuguese: "🇵🇹 Português",
		LangArmenian:   "🇦🇲 Հայերեն",
		LangUkrainian:  "🇺🇦 Українська",
		LangKazakh:     "🇰🇿 Қазақша",
		LangKyrgyz:     "🇰🇬 Кыргызча",
		LangArabic:     "🇸🇦 العربية",
		LangHindi:      "🇮🇳 हिन्दी",
	},
	"it": {
		// Buttons
		ButtonStartConversation: "💬️ Inizia Conversazione",
		ButtonModelSelect:       "🤖 Seleziona Modello",
		ButtonBackToMenu:        "🔙 Torna al Menu",
		ButtonNewConversation:   "➕ Nuova Conversazione",
		ButtonSettings:          "⚙️ Impostazioni",
		ButtonLanguage:          "🌐 Lingua",
		ButtonPrevPage:          "⬅️",
		ButtonNextPage:          "➡️",

		// Menu
		MenuWelcome:    "Benvenuto! Scegli un'opzione:",
		MenuBackToMain: "🏠 Torna al menu principale. Scegli un'opzione:",

		// Conversation
		ConversationStarted:       "🗣️ Conversazione iniziata! Inviami un messaggio e ti risponderò. Puoi sempre tornare al menu.",
		ConversationResumed:       "🗣️ Conversazione ripresa! Inviami un messaggio e ti risponderò. Puoi sempre tornare al menu.",
		ConversationModePrompt:    "Sei in modalità conversazione. Invia un messaggio per chattare, o torna al menu:",
		ConversationNameGenerated: "nome conversazione generato con successo",

		// Settings
		SettingsTitle:         "⚙️ Impostazioni. Scegli un'opzione:",
		LanguageSelectTitle:   "🌐 Seleziona la tua lingua:",
		LanguageUpdateSuccess: "✅ Lingua aggiornata con successo!",

		// Conversation List
		ConversationListSelect:   "💬 Seleziona una conversazione:",
		ConversationListEmpty:    "💬 Nessuna conversazione precedente. Iniziane una nuova:",
		ConversationListPageInfo: "💬 Seleziona una conversazione (pagina %d/%d):",

		// Model Selection
		ModelSelectTitle:   "🤖 Seleziona modello IA:",
		ModelUpdateSuccess: "✅ Modello aggiornato con successo!",

		// Queue
		QueueMessageQueued: "⏳ Il tuo messaggio è stato messo in coda (posizione: %d). Lo elaborerò dopo aver terminato la risposta attuale.",

		// Languages
		LangEnglish:    "🇺🇸 English",
		LangSpanish:    "🇪🇸 Español",
		LangFrench:     "🇫🇷 Français",
		LangGerman:     "🇩🇪 Deutsch",
		LangItalian:    "🇮🇹 Italiano",
		LangRussian:    "🇷🇺 Русский",
		LangChinese:    "🇨🇳 中文",
		LangJapanese:   "🇯🇵 日本語",
		LangKorean:     "🇰🇷 한국어",
		LangPortuguese: "🇵🇹 Português",
		LangArmenian:   "🇦🇲 Հայերեն",
		LangUkrainian:  "🇺🇦 Українська",
		LangKazakh:     "🇰🇿 Қазақша",
		LangKyrgyz:     "🇰🇬 Кыргызча",
		LangArabic:     "🇸🇦 العربية",
		LangHindi:      "🇮🇳 हिन्दी",
	},
	"zh": {
		// Buttons
		ButtonStartConversation: "💬️ 开始对话",
		ButtonModelSelect:       "🤖 选择模型",
		ButtonBackToMenu:        "🔙 返回菜单",
		ButtonNewConversation:   "➕ 新对话",
		ButtonSettings:          "⚙️ 设置",
		ButtonLanguage:          "🌐 语言",
		ButtonPrevPage:          "⬅️",
		ButtonNextPage:          "➡️",

		// Menu
		MenuWelcome:    "欢迎！请选择一个选项：",
		MenuBackToMain: "🏠 返回主菜单。请选择一个选项：",

		// Conversation
		ConversationStarted:       "🗣️ 对话已开始！发送消息，我会回复您。您随时可以返回菜单。",
		ConversationResumed:       "🗣️ 对话已恢复！发送消息，我会回复您。您随时可以返回菜单。",
		ConversationModePrompt:    "您处于对话模式。发送消息进行聊天，或返回菜单：",
		ConversationNameGenerated: "对话名称生成成功",

		// Settings
		SettingsTitle:         "⚙️ 设置。请选择一个选项：",
		LanguageSelectTitle:   "🌐 选择您的语言：",
		LanguageUpdateSuccess: "✅ 语言更新成功！",

		// Conversation List
		ConversationListSelect:   "💬 选择一个对话：",
		ConversationListEmpty:    "💬 无历史对话。开始新对话：",
		ConversationListPageInfo: "💬 选择一个对话（第%d/%d页）：",

		// Model Selection
		ModelSelectTitle:   "🤖 选择AI模型：",
		ModelUpdateSuccess: "✅ 模型更新成功！",

		// Queue
		QueueMessageQueued: "⏳ 您的消息已排队（位置：%d）。我会在完成当前回复后处理它。",

		// Languages
		LangEnglish:    "🇺🇸 English",
		LangSpanish:    "🇪🇸 Español",
		LangFrench:     "🇫🇷 Français",
		LangGerman:     "🇩🇪 Deutsch",
		LangItalian:    "🇮🇹 Italiano",
		LangRussian:    "🇷🇺 Русский",
		LangChinese:    "🇨🇳 中文",
		LangJapanese:   "🇯🇵 日本語",
		LangKorean:     "🇰🇷 한국어",
		LangPortuguese: "🇵🇹 Português",
		LangArmenian:   "🇦🇲 Հայերեն",
		LangUkrainian:  "🇺🇦 Українська",
		LangKazakh:     "🇰🇿 Қазақша",
		LangKyrgyz:     "🇰🇬 Кыргызча",
		LangArabic:     "🇸🇦 العربية",
		LangHindi:      "🇮🇳 हिन्दी",
	},
	"ja": {
		// Buttons
		ButtonStartConversation: "💬️ 会話を始める",
		ButtonModelSelect:       "🤖 モデルを選択",
		ButtonBackToMenu:        "🔙 メニューに戻る",
		ButtonNewConversation:   "➕ 新しい会話",
		ButtonSettings:          "⚙️ 設定",
		ButtonLanguage:          "🌐 言語",
		ButtonPrevPage:          "⬅️",
		ButtonNextPage:          "➡️",

		// Menu
		MenuWelcome:    "ようこそ！オプションを選択してください：",
		MenuBackToMain: "🏠 メインメニューに戻る。オプションを選択してください：",

		// Conversation
		ConversationStarted:       "🗣️ 会話が開始されました！メッセージを送信すると返信します。いつでもメニューに戻れます。",
		ConversationResumed:       "🗣️ 会話が再開されました！メッセージを送信すると返信します。いつでもメニューに戻れます。",
		ConversationModePrompt:    "会話モードです。メッセージを送信してチャットするか、メニューに戻ってください：",
		ConversationNameGenerated: "会話名が正常に生成されました",

		// Settings
		SettingsTitle:         "⚙️ 設定。オプションを選択してください：",
		LanguageSelectTitle:   "🌐 言語を選択してください：",
		LanguageUpdateSuccess: "✅ 言語が正常に更新されました！",

		// Conversation List
		ConversationListSelect:   "💬 会話を選択してください：",
		ConversationListEmpty:    "💬 過去の会話がありません。新しい会話を始めてください：",
		ConversationListPageInfo: "💬 会話を選択してください（%d/%dページ）：",

		// Model Selection
		ModelSelectTitle:   "🤖 AIモデルを選択：",
		ModelUpdateSuccess: "✅ モデルが正常に更新されました！",

		// Queue
		QueueMessageQueued: "⏳ メッセージがキューに追加されました（位置：%d）。現在の応答を完了した後に処理します。",

		// Languages
		LangEnglish:    "🇺🇸 English",
		LangSpanish:    "🇪🇸 Español",
		LangFrench:     "🇫🇷 Français",
		LangGerman:     "🇩🇪 Deutsch",
		LangItalian:    "🇮🇹 Italiano",
		LangRussian:    "🇷🇺 Русский",
		LangChinese:    "🇨🇳 中文",
		LangJapanese:   "🇯🇵 日本語",
		LangKorean:     "🇰🇷 한국어",
		LangPortuguese: "🇵🇹 Português",
		LangArmenian:   "🇦🇲 Հայերեն",
		LangUkrainian:  "🇺🇦 Українська",
		LangKazakh:     "🇰🇿 Қазақша",
		LangKyrgyz:     "🇰🇬 Кыргызча",
		LangArabic:     "🇸🇦 العربية",
		LangHindi:      "🇮🇳 हिन्दी",
	},
	"ko": {
		// Buttons
		ButtonStartConversation: "💬️ 대화 시작",
		ButtonModelSelect:       "🤖 모델 선택",
		ButtonBackToMenu:        "🔙 메뉴로 돌아가기",
		ButtonNewConversation:   "➕ 새 대화",
		ButtonSettings:          "⚙️ 설정",
		ButtonLanguage:          "🌐 언어",
		ButtonPrevPage:          "⬅️",
		ButtonNextPage:          "➡️",

		// Menu
		MenuWelcome:    "환영합니다! 옵션을 선택하세요:",
		MenuBackToMain: "🏠 메인 메뉴로 돌아가기. 옵션을 선택하세요:",

		// Conversation
		ConversationStarted:       "🗣️ 대화가 시작되었습니다! 메시지를 보내시면 답변해드리겠습니다. 언제든지 메뉴로 돌아갈 수 있습니다.",
		ConversationResumed:       "🗣️ 대화가 재개되었습니다! 메시지를 보내시면 답변해드리겠습니다. 언제든지 메뉴로 돌아갈 수 있습니다.",
		ConversationModePrompt:    "대화 모드입니다. 메시지를 보내서 채팅하거나 메뉴로 돌아가세요:",
		ConversationNameGenerated: "대화 이름이 성공적으로 생성되었습니다",

		// Settings
		SettingsTitle:         "⚙️ 설정. 옵션을 선택하세요:",
		LanguageSelectTitle:   "🌐 언어를 선택하세요:",
		LanguageUpdateSuccess: "✅ 언어가 성공적으로 업데이트되었습니다!",

		// Conversation List
		ConversationListSelect:   "💬 대화를 선택하세요:",
		ConversationListEmpty:    "💬 이전 대화가 없습니다. 새 대화를 시작하세요:",
		ConversationListPageInfo: "💬 대화를 선택하세요 (%d/%d페이지):",

		// Model Selection
		ModelSelectTitle:   "🤖 AI 모델 선택:",
		ModelUpdateSuccess: "✅ 모델이 성공적으로 업데이트되었습니다!",

		// Queue
		QueueMessageQueued: "⏳ 메시지가 대기열에 추가되었습니다 (위치: %d). 현재 응답을 완료한 후 처리하겠습니다.",

		// Languages
		LangEnglish:    "🇺🇸 English",
		LangSpanish:    "🇪🇸 Español",
		LangFrench:     "🇫🇷 Français",
		LangGerman:     "🇩🇪 Deutsch",
		LangItalian:    "🇮🇹 Italiano",
		LangRussian:    "🇷🇺 Русский",
		LangChinese:    "🇨🇳 中文",
		LangJapanese:   "🇯🇵 日本語",
		LangKorean:     "🇰🇷 한국어",
		LangPortuguese: "🇵🇹 Português",
		LangArmenian:   "🇦🇲 Հայերեն",
		LangUkrainian:  "🇺🇦 Українська",
		LangKazakh:     "🇰🇿 Қазақша",
		LangKyrgyz:     "🇰🇬 Кыргызча",
		LangArabic:     "🇸🇦 العربية",
		LangHindi:      "🇮🇳 हिन्दी",
	},
	"pt": {
		// Buttons
		ButtonStartConversation: "💬️ Iniciar Conversa",
		ButtonModelSelect:       "🤖 Selecionar Modelo",
		ButtonBackToMenu:        "🔙 Voltar ao Menu",
		ButtonNewConversation:   "➕ Nova Conversa",
		ButtonSettings:          "⚙️ Configurações",
		ButtonLanguage:          "🌐 Idioma",
		ButtonPrevPage:          "⬅️",
		ButtonNextPage:          "➡️",

		// Menu
		MenuWelcome:    "Bem-vindo! Escolha uma opção:",
		MenuBackToMain: "🏠 Voltar ao menu principal. Escolha uma opção:",

		// Conversation
		ConversationStarted:       "🗣️ Conversa iniciada! Envie-me uma mensagem e eu responderei. Você sempre pode voltar ao menu.",
		ConversationResumed:       "🗣️ Conversa retomada! Envie-me uma mensagem e eu responderei. Você sempre pode voltar ao menu.",
		ConversationModePrompt:    "Você está no modo conversa. Envie uma mensagem para conversar, ou volte ao menu:",
		ConversationNameGenerated: "nome da conversa gerado com sucesso",

		// Settings
		SettingsTitle:         "⚙️ Configurações. Escolha uma opção:",
		LanguageSelectTitle:   "🌐 Selecione seu idioma:",
		LanguageUpdateSuccess: "✅ Idioma atualizado com sucesso!",

		// Conversation List
		ConversationListSelect:   "💬 Selecione uma conversa:",
		ConversationListEmpty:    "💬 Nenhuma conversa anterior. Inicie uma nova:",
		ConversationListPageInfo: "💬 Selecione uma conversa (página %d/%d):",

		// Model Selection
		ModelSelectTitle:   "🤖 Selecionar modelo de IA:",
		ModelUpdateSuccess: "✅ Modelo atualizado com sucesso!",

		// Queue
		QueueMessageQueued: "⏳ Sua mensagem foi colocada na fila (posição: %d). Vou processá-la após terminar a resposta atual.",

		// Languages
		LangEnglish:    "🇺🇸 English",
		LangSpanish:    "🇪🇸 Español",
		LangFrench:     "🇫🇷 Français",
		LangGerman:     "🇩🇪 Deutsch",
		LangItalian:    "🇮🇹 Italiano",
		LangRussian:    "🇷🇺 Русский",
		LangChinese:    "🇨🇳 中文",
		LangJapanese:   "🇯🇵 日本語",
		LangKorean:     "🇰🇷 한국어",
		LangPortuguese: "🇵🇹 Português",
		LangArmenian:   "🇦🇲 Հայերեն",
		LangUkrainian:  "🇺🇦 Українська",
		LangKazakh:     "🇰🇿 Қазақша",
		LangKyrgyz:     "🇰🇬 Кыргызча",
		LangArabic:     "🇸🇦 العربية",
		LangHindi:      "🇮🇳 हिन्दी",
	},
	"hy": {
		// Buttons
		ButtonStartConversation: "💬️ Սկսել խոսակցություն",
		ButtonModelSelect:       "🤖 Ընտրել մոդելը",
		ButtonBackToMenu:        "🔙 Վերադառնալ մենյու",
		ButtonNewConversation:   "➕ Նոր խոսակցություն",
		ButtonSettings:          "⚙️ Կարգավորումներ",
		ButtonLanguage:          "🌐 Լեզու",
		ButtonPrevPage:          "⬅️",
		ButtonNextPage:          "➡️",

		// Menu
		MenuWelcome:    "Բարի գալուստ: Ընտրեք տարբերակը:",
		MenuBackToMain: "🏠 Վերադառնալ գլխավոր մենյու: Ընտրեք տարբերակը:",

		// Conversation
		ConversationStarted:       "🗣️ Խոսակցությունը սկսվեց: Ուղարկեք ինձ հաղորդագրություն և ես կպատասխանեմ: Դուք միշտ կարող եք վերադառնալ մենյու:",
		ConversationResumed:       "🗣️ Խոսակցությունը շարունակվեց: Ուղարկեք ինձ հաղորդագրություն և ես կպատասխանեմ: Դուք միշտ կարող եք վերադառնալ մենյու:",
		ConversationModePrompt:    "Դուք խոսակցության ռեժիմում եք: Ուղարկեք հաղորդագրություն՝ չատ անելու համար, կամ վերադառնալ մենյու:",
		ConversationNameGenerated: "խոսակցության անունը հաջողությամբ ստեղծվեց",

		// Settings
		SettingsTitle:         "⚙️ Կարգավորումներ: Ընտրեք տարբերակը:",
		LanguageSelectTitle:   "🌐 Ընտրեք ձեր լեզուն:",
		LanguageUpdateSuccess: "✅ Լեզուն հաջողությամբ թարմացվեց:",

		// Conversation List
		ConversationListSelect:   "💬 Ընտրեք խոսակցությունը:",
		ConversationListEmpty:    "💬 Նախկին խոսակցություններ չկան: Սկսեք նորը:",
		ConversationListPageInfo: "💬 Ընտրեք խոսակցությունը (էջ %d/%d):",

		// Model Selection
		ModelSelectTitle:   "🤖 Ընտրել AI մոդելը:",
		ModelUpdateSuccess: "✅ Մոդելը հաջողությամբ թարմացվեց:",

		// Queue
		QueueMessageQueued: "⏳ Ձեր հաղորդագրությունը հերթի մեջ է (դիրքը՝ %d): Ես կմշակեմ այն ընթացիկ պատասխանն ավարտելուց հետո:",

		// Languages
		LangEnglish:    "🇺🇸 English",
		LangSpanish:    "🇪🇸 Español",
		LangFrench:     "🇫🇷 Français",
		LangGerman:     "🇩🇪 Deutsch",
		LangItalian:    "🇮🇹 Italiano",
		LangRussian:    "🇷🇺 Русский",
		LangChinese:    "🇨🇳 中文",
		LangJapanese:   "🇯🇵 日本語",
		LangKorean:     "🇰🇷 한국어",
		LangPortuguese: "🇵🇹 Português",
		LangArmenian:   "🇦🇲 Հայերեն",
		LangUkrainian:  "🇺🇦 Українська",
		LangKazakh:     "🇰🇿 Қазақша",
		LangKyrgyz:     "🇰🇬 Кыргызча",
		LangArabic:     "🇸🇦 العربية",
		LangHindi:      "🇮🇳 हिन्दी",
	},
	"uk": {
		// Buttons
		ButtonStartConversation: "💬️ Почати розмову",
		ButtonModelSelect:       "🤖 Вибрати модель",
		ButtonBackToMenu:        "🔙 Назад до меню",
		ButtonNewConversation:   "➕ Нова розмова",
		ButtonSettings:          "⚙️ Налаштування",
		ButtonLanguage:          "🌐 Мова",
		ButtonPrevPage:          "⬅️",
		ButtonNextPage:          "➡️",

		// Menu
		MenuWelcome:    "Вітаємо! Оберіть опцію:",
		MenuBackToMain: "🏠 Назад до головного меню. Оберіть опцію:",

		// Conversation
		ConversationStarted:       "🗣️ Розмову розпочато! Надішліть мені повідомлення, і я відповім. Ви завжди можете повернутися до меню.",
		ConversationResumed:       "🗣️ Розмову відновлено! Надішліть мені повідомлення, і я відповім. Ви завжди можете повернутися до меню.",
		ConversationModePrompt:    "Ви в режимі розмови. Надішліть повідомлення для чату або поверніться до меню:",
		ConversationNameGenerated: "назву розмови успішно створено",

		// Settings
		SettingsTitle:         "⚙️ Налаштування. Оберіть опцію:",
		LanguageSelectTitle:   "🌐 Оберіть вашу мову:",
		LanguageUpdateSuccess: "✅ Мову успішно оновлено!",

		// Conversation List
		ConversationListSelect:   "💬 Оберіть розмову:",
		ConversationListEmpty:    "💬 Немає попередніх розмов. Почніть нову:",
		ConversationListPageInfo: "💬 Оберіть розмову (сторінка %d/%d):",

		// Model Selection
		ModelSelectTitle:   "🤖 Вибрати AI модель:",
		ModelUpdateSuccess: "✅ Модель успішно оновлено!",

		// Queue
		QueueMessageQueued: "⏳ Ваше повідомлення поставлено в чергу (позиція: %d). Я оброблю його після завершення поточної відповіді.",

		// Languages
		LangEnglish:    "🇺🇸 English",
		LangSpanish:    "🇪🇸 Español",
		LangFrench:     "🇫🇷 Français",
		LangGerman:     "🇩🇪 Deutsch",
		LangItalian:    "🇮🇹 Italiano",
		LangRussian:    "🇷🇺 Русский",
		LangChinese:    "🇨🇳 中文",
		LangJapanese:   "🇯🇵 日本語",
		LangKorean:     "🇰🇷 한국어",
		LangPortuguese: "🇵🇹 Português",
		LangArmenian:   "🇦🇲 Հայերեն",
		LangUkrainian:  "🇺🇦 Українська",
		LangKazakh:     "🇰🇿 Қазақша",
		LangKyrgyz:     "🇰🇬 Кыргызча",
		LangArabic:     "🇸🇦 العربية",
		LangHindi:      "🇮🇳 हिन्दी",
	},
	"kk": {
		// Buttons
		ButtonStartConversation: "💬️ Сөйлесуді бастау",
		ButtonModelSelect:       "🤖 Модельді таңдау",
		ButtonBackToMenu:        "🔙 Мәзірге оралу",
		ButtonNewConversation:   "➕ Жаңа сөйлесу",
		ButtonSettings:          "⚙️ Баптаулар",
		ButtonLanguage:          "🌐 Тіл",
		ButtonPrevPage:          "⬅️",
		ButtonNextPage:          "➡️",

		// Menu
		MenuWelcome:    "Қош келдіңіз! Опцияны таңдаңыз:",
		MenuBackToMain: "🏠 Басты мәзірге оралу. Опцияны таңдаңыз:",

		// Conversation
		ConversationStarted:       "🗣️ Сөйлесу басталды! Маған хабарлама жіберіңіз, мен жауап беремін. Сіз әрқашан мәзірге орала аласыз.",
		ConversationResumed:       "🗣️ Сөйлесу жалғасты! Маған хабарлама жіберіңіз, мен жауап беремін. Сіз әрқашан мәзірге орала аласыз.",
		ConversationModePrompt:    "Сіз сөйлесу режимінде жүрсіз. Чат үшін хабарлама жіберіңіз немесе мәзірге оралыңыз:",
		ConversationNameGenerated: "сөйлесу атауы сәтті құрылды",

		// Settings
		SettingsTitle:         "⚙️ Баптаулар. Опцияны таңдаңыз:",
		LanguageSelectTitle:   "🌐 Тіліңізді таңдаңыз:",
		LanguageUpdateSuccess: "✅ Тіл сәтті жаңартылды!",

		// Conversation List
		ConversationListSelect:   "💬 Сөйлесуді таңдаңыз:",
		ConversationListEmpty:    "💬 Бұрынғы сөйлесулер жоқ. Жаңасын бастаңыз:",
		ConversationListPageInfo: "💬 Сөйлесуді таңдаңыз (%d/%d бет):",

		// Model Selection
		ModelSelectTitle:   "🤖 AI моделін таңдау:",
		ModelUpdateSuccess: "✅ Модель сәтті жаңартылды!",

		// Queue
		QueueMessageQueued: "⏳ Сіздің хабарламаңыз кезекке қойылды (орын: %d). Мен оны ағымдағы жауапты аяқтағаннан кейін өңдеймін.",

		// Languages
		LangEnglish:    "🇺🇸 English",
		LangSpanish:    "🇪🇸 Español",
		LangFrench:     "🇫🇷 Français",
		LangGerman:     "🇩🇪 Deutsch",
		LangItalian:    "🇮🇹 Italiano",
		LangRussian:    "🇷🇺 Русский",
		LangChinese:    "🇨🇳 中文",
		LangJapanese:   "🇯🇵 日本語",
		LangKorean:     "🇰🇷 한국어",
		LangPortuguese: "🇵🇹 Português",
		LangArmenian:   "🇦🇲 Հայերեն",
		LangUkrainian:  "🇺🇦 Українська",
		LangKazakh:     "🇰🇿 Қазақша",
		LangKyrgyz:     "🇰🇬 Кыргызча",
		LangArabic:     "🇸🇦 العربية",
		LangHindi:      "🇮🇳 हिन्दी",
	},
	"ky": {
		// Buttons
		ButtonStartConversation: "💬️ Маекти баштоо",
		ButtonModelSelect:       "🤖 Модель тандоо",
		ButtonBackToMenu:        "🔙 Менюга кайтуу",
		ButtonNewConversation:   "➕ Жаңы маек",
		ButtonSettings:          "⚙️ Жөндөөлөр",
		ButtonLanguage:          "🌐 Тил",
		ButtonPrevPage:          "⬅️",
		ButtonNextPage:          "➡️",

		// Menu
		MenuWelcome:    "Кош келиңиз! Опцияны тандаңыз:",
		MenuBackToMain: "🏠 Башкы менюга кайтуу. Опцияны тандаңыз:",

		// Conversation
		ConversationStarted:       "🗣️ Маек башталды! Мага билдирүү жөнөтүңүз, мен жооп берем. Сиз дайыма менюга кайта аласыз.",
		ConversationResumed:       "🗣️ Маек улантылды! Мага билдирүү жөнөтүңүз, мен жооп берем. Сиз дайыма менюга кайта аласыз.",
		ConversationModePrompt:    "Сиз маек режиминдесиз. Сүйлөшүү үчүн билдирүү жөнөтүңүз же менюга кайтыңыз:",
		ConversationNameGenerated: "маек аталышы ийгиликтүү түзүлдү",

		// Settings
		SettingsTitle:         "⚙️ Жөндөөлөр. Опцияны тандаңыз:",
		LanguageSelectTitle:   "🌐 Тилиңизди тандаңыз:",
		LanguageUpdateSuccess: "✅ Тил ийгиликтүү жаңыртылды!",

		// Conversation List
		ConversationListSelect:   "💬 Маекти тандаңыз:",
		ConversationListEmpty:    "💬 Мурунку маектер жок. Жаңысын баштаңыз:",
		ConversationListPageInfo: "💬 Маекти тандаңыз (%d/%d бет):",

		// Model Selection
		ModelSelectTitle:   "🤖 AI модель тандоо:",
		ModelUpdateSuccess: "✅ Модель ийгиликтүү жаңыртылды!",

		// Queue
		QueueMessageQueued: "⏳ Сиздин билдирүүңүз кезекке коюлду (орун: %d). Мен аны учурдагы жоопту бүткөндөн кийин иштетем.",

		// Languages
		LangEnglish:    "🇺🇸 English",
		LangSpanish:    "🇪🇸 Español",
		LangFrench:     "🇫🇷 Français",
		LangGerman:     "🇩🇪 Deutsch",
		LangItalian:    "🇮🇹 Italiano",
		LangRussian:    "🇷🇺 Русский",
		LangChinese:    "🇨🇳 中文",
		LangJapanese:   "🇯🇵 日本語",
		LangKorean:     "🇰🇷 한국어",
		LangPortuguese: "🇵🇹 Português",
		LangArmenian:   "🇦🇲 Հայերեն",
		LangUkrainian:  "🇺🇦 Українська",
		LangKazakh:     "🇰🇿 Қазақша",
		LangKyrgyz:     "🇰🇬 Кыргызча",
		LangArabic:     "🇸🇦 العربية",
		LangHindi:      "🇮🇳 हिन्दी",
	},
	"ar": {
		// Buttons
		ButtonStartConversation: "💬️ بدء محادثة",
		ButtonModelSelect:       "🤖 اختيار نموذج",
		ButtonBackToMenu:        "🔙 العودة للقائمة",
		ButtonNewConversation:   "➕ محادثة جديدة",
		ButtonSettings:          "⚙️ الإعدادات",
		ButtonLanguage:          "🌐 اللغة",
		ButtonPrevPage:          "⬅️",
		ButtonNextPage:          "➡️",

		// Menu
		MenuWelcome:    "مرحباً! اختر خياراً:",
		MenuBackToMain: "🏠 العودة للقائمة الرئيسية. اختر خياراً:",

		// Conversation
		ConversationStarted:       "🗣️ بدأت المحادثة! أرسل لي رسالة وسأرد عليك. يمكنك العودة للقائمة في أي وقت.",
		ConversationResumed:       "🗣️ استُأنِفت المحادثة! أرسل لي رسالة وسأرد عليك. يمكنك العودة للقائمة في أي وقت.",
		ConversationModePrompt:    "أنت في وضع المحادثة. أرسل رسالة للدردشة أو ارجع للقائمة:",
		ConversationNameGenerated: "تم إنشاء اسم المحادثة بنجاح",

		// Settings
		SettingsTitle:         "⚙️ الإعدادات. اختر خياراً:",
		LanguageSelectTitle:   "🌐 اختر لغتك:",
		LanguageUpdateSuccess: "✅ تم تحديث اللغة بنجاح!",

		// Conversation List
		ConversationListSelect:   "💬 اختر محادثة:",
		ConversationListEmpty:    "💬 لا توجد محادثات سابقة. ابدأ محادثة جديدة:",
		ConversationListPageInfo: "💬 اختر محادثة (الصفحة %d/%d):",

		// Model Selection
		ModelSelectTitle:   "🤖 اختيار نموذج الذكاء الاصطناعي:",
		ModelUpdateSuccess: "✅ تم تحديث النموذج بنجاح!",

		// Queue
		QueueMessageQueued: "⏳ تم وضع رسالتك في طابور الانتظار (الموضع: %d). سأقوم بمعالجتها بعد إنهاء الرد الحالي.",

		// Languages
		LangEnglish:    "🇺🇸 English",
		LangSpanish:    "🇪🇸 Español",
		LangFrench:     "🇫🇷 Français",
		LangGerman:     "🇩🇪 Deutsch",
		LangItalian:    "🇮🇹 Italiano",
		LangRussian:    "🇷🇺 Русский",
		LangChinese:    "🇨🇳 中文",
		LangJapanese:   "🇯🇵 日本語",
		LangKorean:     "🇰🇷 한국어",
		LangPortuguese: "🇵🇹 Português",
		LangArmenian:   "🇦🇲 Հայերեն",
		LangUkrainian:  "🇺🇦 Українська",
		LangKazakh:     "🇰🇿 Қазақша",
		LangKyrgyz:     "🇰🇬 Кыргызча",
		LangArabic:     "🇸🇦 العربية",
		LangHindi:      "🇮🇳 हिन्दी",
	},
	"hi": {
		// Buttons
		ButtonStartConversation: "💬️ बातचीत शुरू करें",
		ButtonModelSelect:       "🤖 मॉडल चुनें",
		ButtonBackToMenu:        "🔙 मेनू में वापस जाएं",
		ButtonNewConversation:   "➕ नई बातचीत",
		ButtonSettings:          "⚙️ सेटिंग्स",
		ButtonLanguage:          "🌐 भाषा",
		ButtonPrevPage:          "⬅️",
		ButtonNextPage:          "➡️",

		// Menu
		MenuWelcome:    "स्वागत है! एक विकल्प चुनें:",
		MenuBackToMain: "🏠 मुख्य मेनू में वापस जाएं। एक विकल्प चुनें:",

		// Conversation
		ConversationStarted:       "🗣️ बातचीत शुरू हो गई! मुझे संदेश भेजें और मैं जवाब दूंगा। आप कभी भी मेनू में वापस जा सकते हैं।",
		ConversationResumed:       "🗣️ बातचीत फिर से शुरू! मुझे संदेश भेजें और मैं जवाब दूंगा। आप कभी भी मेनू में वापस जा सकते हैं।",
		ConversationModePrompt:    "आप बातचीत मोड में हैं। चैट करने के लिए संदेश भेजें, या मेनू में वापस जाएं:",
		ConversationNameGenerated: "बातचीत का नाम सफलतापूर्वक बनाया गया",

		// Settings
		SettingsTitle:         "⚙️ सेटिंग्स। एक विकल्प चुनें:",
		LanguageSelectTitle:   "🌐 अपनी भाषा चुनें:",
		LanguageUpdateSuccess: "✅ भाषा सफलतापूर्वक अपडेट हो गई!",

		// Conversation List
		ConversationListSelect:   "💬 एक बातचीत चुनें:",
		ConversationListEmpty:    "💬 कोई पिछली बातचीत नहीं। नई शुरू करें:",
		ConversationListPageInfo: "💬 एक बातचीत चुनें (पेज %d/%d):",

		// Model Selection
		ModelSelectTitle:   "🤖 AI मॉडल चुनें:",
		ModelUpdateSuccess: "✅ मॉडल सफलतापूर्वक अपडेट हो गया!",

		// Queue
		QueueMessageQueued: "⏳ आपका संदेश कतार में रखा गया है (स्थिति: %d)। मैं वर्तमान जवाब पूरा करने के बाद इसे प्रोसेस करूंगा।",

		// Languages
		LangEnglish:    "🇺🇸 English",
		LangSpanish:    "🇪🇸 Español",
		LangFrench:     "🇫🇷 Français",
		LangGerman:     "🇩🇪 Deutsch",
		LangItalian:    "🇮🇹 Italiano",
		LangRussian:    "🇷🇺 Русский",
		LangChinese:    "🇹🇼 中文",
		LangJapanese:   "🇯🇵 日本語",
		LangKorean:     "🇰🇷 한국어",
		LangPortuguese: "🇵🇹 Português",
		LangArmenian:   "🇦🇲 Հայերեն",
		LangUkrainian:  "🇺🇦 Українська",
		LangKazakh:     "🇰🇿 Қазақша",
		LangKyrgyz:     "🇰🇬 Кыргызча",
		LangArabic:     "🇸🇦 العربية",
		LangHindi:      "🇮🇳 हिन्दी",
	},
}

// GetString returns a localized string for the given key and language.
// Falls back to English if the language or key is not found.
func GetString(lang, key string) string {
	if langStrings, exists := Strings[lang]; exists {
		if str, found := langStrings[key]; found {
			return str
		}
	}

	// Fallback to English
	if engStrings, exists := Strings["en"]; exists {
		if str, found := engStrings[key]; found {
			return str
		}
	}

	// Last resort fallback
	return key
}

// GetLanguageDisplayName returns the display name for a language key.
func GetLanguageDisplayName(lang, langKey string) string {
	return GetString(lang, langKey)
}
