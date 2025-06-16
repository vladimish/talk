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
	ButtonProfile           = "button.profile"
	ButtonSubscription      = "button.subscription"
	ButtonPrevPage          = "button.prev_page"
	ButtonNextPage          = "button.next_page"

	// Menu messages.
	MenuWelcome    = "menu.welcome"
	MenuBackToMain = "menu.back_to_main"

	// Conversation messages.
	ConversationStarted             = "conversation.started"
	ConversationResumed             = "conversation.resumed"
	ConversationModePrompt          = "conversation.mode_prompt"
	ConversationNameGenerated       = "conversation.name_generated"
	ConversationStartedWebSearchOff = "conversation.started_web_search_off"
	ConversationStartedWebSearchOn  = "conversation.started_web_search_on"
	ConversationResumedWebSearchOff = "conversation.resumed_web_search_off"
	ConversationResumedWebSearchOn  = "conversation.resumed_web_search_on"

	// Settings messages.
	SettingsTitle         = "settings.title"
	LanguageSelectTitle   = "settings.language_select_title"
	LanguageUpdateSuccess = "settings.language_update_success"

	// Conversation list messages.
	ConversationListSelect   = "conversation_list.select"
	ConversationListEmpty    = "conversation_list.empty"
	ConversationListPageInfo = "conversation_list.page_info"

	// Model selection messages.
	ModelSelectTitle       = "model.select_title"
	ModelUpdateSuccess     = "model.update_success"
	ModelImageNotSupported = "model.image_not_supported"
	ModelPDFNotSupported   = "model.pdf_not_supported"

	// Model names (for internationalization).
	ModelGeminiFlash   = "model.gemini_flash"
	ModelGPT4o         = "model.gpt4o"
	ModelGPT4oMini     = "model.gpt4o_mini"
	ModelClaude4Sonnet = "model.claude4_sonnet"
	ModelGeminiPro     = "model.gemini_pro"
	ModelO3Mini        = "model.o3_mini"
	ModelDeepSeekChat  = "model.deepseek_chat"
	ModelDeepSeekR1    = "model.deepseek_r1"

	// Queue messages.
	QueueMessageQueued = "queue.message_queued"

	// Web search messages.
	ButtonWebSearchOn             = "button.web_search_on"
	ButtonWebSearchOff            = "button.web_search_off"
	WebSearchSubscriptionRequired = "web_search.subscription_required"
	WebSearchEnabledNotification  = "web_search.enabled_notification"

	// Profile messages.
	ProfileTitle              = "profile.title"
	ProfileTokenBalance       = "profile.token_balance"
	ProfilePremiumTokens      = "profile.premium_tokens"
	ProfileRegularTokens      = "profile.regular_tokens"
	ProfileInsufficientTokens = "profile.insufficient_tokens" //nolint:gosec

	// Subscription messages.
	SubscriptionTitle        = "subscription.title"
	SubscriptionMonthlyOffer = "subscription.monthly_offer"
	SubscriptionBuyButton    = "subscription.buy_button"
	SubscriptionSuccess      = "subscription.success"
	SubscriptionFailed       = "subscription.failed"
	SubscriptionActiveInfo   = "subscription.active_info"
	SubscriptionExpired      = "subscription.expired"

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
		ButtonProfile:           "👤 Profile",
		ButtonSubscription:      "💳 Subscription",
		ButtonPrevPage:          "⬅️",
		ButtonNextPage:          "➡️",

		// Menu
		MenuWelcome:    "Welcome! Choose an option:",
		MenuBackToMain: "🏠 Back to main menu. Choose an option:",

		// Conversation
		ConversationStarted:             "🗣️ Conversation started! Send me a message and I'll respond. You can always go back to the menu.",
		ConversationResumed:             "🗣️ Conversation resumed! Send me a message and I'll respond. You can always go back to the menu.",
		ConversationModePrompt:          "You're in conversation mode. Send a message to chat, or go back to menu:",
		ConversationNameGenerated:       "conversation name generated successfully",
		ConversationStartedWebSearchOff: "🗣️ Conversation started! Send me a message and I'll respond. Web search is OFF.",
		ConversationStartedWebSearchOn:  "🗣️ Conversation started! Send me a message and I'll respond. Web search is ON (costs 1 premium token per query).",
		ConversationResumedWebSearchOff: "🗣️ Conversation resumed! Send me a message and I'll respond. Web search is OFF.",
		ConversationResumedWebSearchOn:  "🗣️ Conversation resumed! Send me a message and I'll respond. Web search is ON (costs 1 premium token per query).",

		// Settings
		SettingsTitle:         "⚙️ Settings. Choose an option:",
		LanguageSelectTitle:   "🌐 Select your language:",
		LanguageUpdateSuccess: "✅ Language updated successfully!",

		// Conversation List
		ConversationListSelect:   "💬 Select a conversation:",
		ConversationListEmpty:    "💬 No previous conversations. Start a new one:",
		ConversationListPageInfo: "💬 Select a conversation (page %d/%d):",

		// Model Selection
		ModelSelectTitle:       "🤖 Select AI model:\n\n👁️ = Image support | 📄 = PDF support | 🧠 = Reasoning | 🌐 = Web search\n❌ = Not available (insufficient tokens/subscription required)",
		ModelUpdateSuccess:     "✅ Model updated successfully!",
		ModelImageNotSupported: "❌ The selected model does not support image inputs. Please choose a different model or send a text message.",
		ModelPDFNotSupported:   "❌ The selected model does not support PDF inputs. Please choose a different model or send a text message.",

		// Model Names
		ModelGeminiFlash:   "🚀 Gemini 2.5 Flash (Fast & efficient for quick responses)",
		ModelGPT4o:         "🧠 GPT-4o (Most capable for complex tasks)",
		ModelGPT4oMini:     "⚡ GPT-4o Mini (Balanced speed & performance)",
		ModelClaude4Sonnet: "🎭 Claude 4 Sonnet (Creative writing & analysis)",
		ModelGeminiPro:     "🌸 Gemini 2.5 Pro (Long documents & context)",
		ModelO3Mini:        "🤖 OpenAI o3-mini (Advanced reasoning model)",
		ModelDeepSeekChat:  "🔬 DeepSeek Chat v3 (Research & coding)",
		ModelDeepSeekR1:    "🔍 DeepSeek R1 (Deep reasoning & logic)",

		// Queue
		QueueMessageQueued: "⏳ Your message has been queued (position: %d). I'll process it after finishing the current response.",

		// Web Search
		ButtonWebSearchOn:             "🌐 Web search: ON",
		ButtonWebSearchOff:            "🌐 Web search: OFF",
		WebSearchSubscriptionRequired: "🔐 Web search requires an active subscription. Please subscribe to use this feature.",
		WebSearchEnabledNotification:  "🌐 Web search enabled! Each search query will cost 1 additional premium token.",

		// Profile
		ProfileTitle:              "👤 Your Profile",
		ProfileTokenBalance:       "💰 Token Balance:",
		ProfilePremiumTokens:      "🟡 Premium: %d tokens",
		ProfileRegularTokens:      "🔵 Regular: %d tokens",
		ProfileInsufficientTokens: "❌ Insufficient tokens. You need %d %s tokens to use this model.",

		// Subscription
		SubscriptionTitle:        "💳 Subscription",
		SubscriptionMonthlyOffer: "🌟 Monthly Premium Subscription\n\n✨ Get 1500 regular tokens + 100 premium tokens every month!\n\nPrice: ⭐ 600 Telegram Stars per month",
		SubscriptionBuyButton:    "💰 Subscribe for ⭐ 600 Stars",
		SubscriptionSuccess:      "🎉 Subscription activated! You've received 1500 regular tokens and 100 premium tokens.",
		SubscriptionFailed:       "❌ Subscription failed. Please try again.",
		SubscriptionActiveInfo:   "✅ You have an active subscription! %d days remaining.",
		SubscriptionExpired:      "❌ Your subscription has expired. Subscribe again to continue receiving tokens.",

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
		ButtonProfile:           "👤 Perfil",
		ButtonSubscription:      "💳 Suscripción",
		ButtonPrevPage:          "⬅️",
		ButtonNextPage:          "➡️",

		// Menu
		MenuWelcome:    "¡Bienvenido! Elige una opción:",
		MenuBackToMain: "🏠 Volver al menú principal. Elige una opción:",

		// Conversation
		ConversationStarted:             "🗣️ ¡Conversación iniciada! Envíame un mensaje y te responderé. Siempre puedes volver al menú.",
		ConversationResumed:             "🗣️ ¡Conversación reanudada! Envíame un mensaje y te responderé. Siempre puedes volver al menú.",
		ConversationModePrompt:          "Estás en modo conversación. Envía un mensaje para chatear, o vuelve al menú:",
		ConversationNameGenerated:       "nombre de conversación generado exitosamente",
		ConversationStartedWebSearchOff: "🗣️ ¡Conversación iniciada! Envíame un mensaje y te responderé. Búsqueda web está DESACTIVADA.",
		ConversationStartedWebSearchOn:  "🗣️ ¡Conversación iniciada! Envíame un mensaje y te responderé. Búsqueda web está ACTIVADA (cuesta 1 token premium por consulta).",
		ConversationResumedWebSearchOff: "🗣️ ¡Conversación reanudada! Envíame un mensaje y te responderé. Búsqueda web está DESACTIVADA.",
		ConversationResumedWebSearchOn:  "🗣️ ¡Conversación reanudada! Envíame un mensaje y te responderé. Búsqueda web está ACTIVADA (cuesta 1 token premium por consulta).",

		// Settings
		SettingsTitle:         "⚙️ Configuración. Elige una opción:",
		LanguageSelectTitle:   "🌐 Selecciona tu idioma:",
		LanguageUpdateSuccess: "✅ ¡Idioma actualizado exitosamente!",

		// Conversation List
		ConversationListSelect:   "💬 Selecciona una conversación:",
		ConversationListEmpty:    "💬 No hay conversaciones anteriores. Inicia una nueva:",
		ConversationListPageInfo: "💬 Selecciona una conversación (página %d/%d):",

		// Model Selection
		ModelSelectTitle:       "🤖 Seleccionar modelo de IA:\n\n👁️ = Soporte de imágenes | 📄 = Soporte de PDF | 🧠 = Razonamiento | 🌐 = Búsqueda web\n❌ = No disponible (tokens insuficientes/suscripción requerida)",
		ModelUpdateSuccess:     "✅ ¡Modelo actualizado exitosamente!",
		ModelImageNotSupported: "❌ El modelo seleccionado no admite imágenes. Por favor elige un modelo diferente o envía un mensaje de texto.",
		ModelPDFNotSupported:   "❌ El modelo seleccionado no admite archivos PDF. Por favor elige un modelo diferente o envía un mensaje de texto.",

		// Model Names
		ModelGeminiFlash:   "🚀 Gemini 2.5 Flash (Rápido y eficiente para respuestas ágiles)",
		ModelGPT4o:         "🧠 GPT-4o (Más capaz para tareas complejas)",
		ModelGPT4oMini:     "⚡ GPT-4o Mini (Velocidad y rendimiento equilibrados)",
		ModelClaude4Sonnet: "🎭 Claude 4 Sonnet (Escritura creativa y análisis)",
		ModelGeminiPro:     "🌸 Gemini 2.5 Pro (Documentos largos y contexto)",
		ModelO3Mini:        "🤖 OpenAI o3-mini (Modelo de razonamiento avanzado)",
		ModelDeepSeekChat:  "🔬 DeepSeek Chat v3 (Investigación y programación)",
		ModelDeepSeekR1:    "🔍 DeepSeek R1 (Razonamiento profundo y lógica)",

		// Queue
		QueueMessageQueued: "⏳ Tu mensaje ha sido puesto en cola (posición: %d). Lo procesaré después de terminar la respuesta actual.",

		// Web Search
		ButtonWebSearchOn:             "🌐 Búsqueda web: ACTIVADA",
		ButtonWebSearchOff:            "🌐 Búsqueda web: DESACTIVADA",
		WebSearchSubscriptionRequired: "🔐 La búsqueda web requiere una suscripción activa. Por favor suscríbete para usar esta característica.",
		WebSearchEnabledNotification:  "🌐 Web search enabled! Each search query will cost 1 additional premium token.",

		// Profile
		ProfileTitle:              "👤 Tu Perfil",
		ProfileTokenBalance:       "💰 Balance de Tokens:",
		ProfilePremiumTokens:      "🟡 Premium: %d tokens",
		ProfileRegularTokens:      "🔵 Regular: %d tokens",
		ProfileInsufficientTokens: "❌ Tokens insuficientes. Necesitas %d tokens %s para usar este modelo.",

		// Subscription
		SubscriptionTitle:        "💳 Suscripción",
		SubscriptionMonthlyOffer: "🌟 Suscripción Premium Mensual\n\n✨ ¡Obtén 1500 tokens regulares + 100 tokens premium cada mes!\n\nPrecio: ⭐ 600 Estrellas de Telegram por mes",
		SubscriptionBuyButton:    "💰 Suscribirse por ⭐ 600 Estrellas",
		SubscriptionSuccess:      "🎉 ¡Suscripción activada! Has recibido 1500 tokens regulares y 100 tokens premium.",
		SubscriptionFailed:       "❌ La suscripción falló. Por favor, inténtalo de nuevo.",
		SubscriptionActiveInfo:   "✅ ¡Tienes una suscripción activa! %d días restantes.",
		SubscriptionExpired:      "❌ Tu suscripción ha expirado. Suscríbete de nuevo para seguir recibiendo tokens.",

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
		ButtonProfile:           "👤 Профиль",
		ButtonSubscription:      "💳 Подписка",
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
		ModelSelectTitle:       "🤖 Выбрать ИИ модель:\n\n👁️ = Поддержка изображений | 📄 = Поддержка PDF | 🧠 = Рассуждение | 🌐 = Веб-поиск\n❌ = Недоступно (недостаточно токенов/требуется подписка)",
		ModelUpdateSuccess:     "✅ Модель успешно обновлена!",
		ModelImageNotSupported: "❌ Выбранная модель не поддерживает изображения. Пожалуйста, выберите другую модель или отправьте текстовое сообщение.",
		ModelPDFNotSupported:   "❌ Выбранная модель не поддерживает PDF файлы. Пожалуйста, выберите другую модель или отправьте текстовое сообщение.",

		// Model Names
		ModelGeminiFlash:   "🚀 Gemini 2.5 Flash (Быстрый и эффективный для оперативных ответов)",
		ModelGPT4o:         "🧠 GPT-4o (Самый способный для сложных задач)",
		ModelGPT4oMini:     "⚡ GPT-4o Mini (Сбалансированная скорость и производительность)",
		ModelClaude4Sonnet: "🎭 Claude 4 Sonnet (Творческое письмо и анализ)",
		ModelGeminiPro:     "🌸 Gemini 2.5 Pro (Длинные документы и контекст)",
		ModelO3Mini:        "🤖 OpenAI o3-mini (Продвинутая модель рассуждений)",
		ModelDeepSeekChat:  "🔬 DeepSeek Chat v3 (Исследования и программирование)",
		ModelDeepSeekR1:    "🔍 DeepSeek R1 (Глубокие рассуждения и логика)",

		// Queue
		QueueMessageQueued: "⏳ Ваше сообщение поставлено в очередь (позиция: %d). Я обработаю его после завершения текущего ответа.",

		// Web Search
		ButtonWebSearchOn:             "🌐 Веб-поиск: ВКЛЮЧЕН",
		ButtonWebSearchOff:            "🌐 Веб-поиск: ВЫКЛЮЧЕН",
		WebSearchSubscriptionRequired: "🔐 Веб-поиск требует активной подписки. Пожалуйста, оформите подписку для использования этой функции.",
		WebSearchEnabledNotification:  "🌐 Web search enabled! Each search query will cost 1 additional premium token.",

		// Profile
		ProfileTitle:              "👤 Ваш Профиль",
		ProfileTokenBalance:       "💰 Баланс Токенов:",
		ProfilePremiumTokens:      "🟡 Премиум: %d токенов",
		ProfileRegularTokens:      "🔵 Обычные: %d токенов",
		ProfileInsufficientTokens: "❌ Недостаточно токенов. Вам нужно %d %s токенов для использования этой модели.",

		// Subscription
		SubscriptionTitle:        "💳 Подписка",
		SubscriptionMonthlyOffer: "🌟 Ежемесячная Премиум Подписка\n\n✨ Получайте 1500 обычных токенов + 100 премиум токенов каждый месяц!\n\nЦена: ⭐ 600 Звезды Telegram в месяц",
		SubscriptionBuyButton:    "💰 Подписаться за ⭐ 600 Звезды",
		SubscriptionSuccess:      "🎉 Подписка активирована! Вы получили 1500 обычных токенов и 100 премиум токенов.",
		SubscriptionFailed:       "❌ Подписка не удалась. Пожалуйста, попробуйте снова.",
		SubscriptionActiveInfo:   "✅ У вас активная подписка! Осталось %d дней.",
		SubscriptionExpired:      "❌ Ваша подписка истекла. Подпишитесь снова, чтобы продолжить получать токены.",

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
		ButtonProfile:           "👤 Profil",
		ButtonSubscription:      "💳 Abonnement",
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
		ModelSelectTitle:       "🤖 Sélectionner un modèle IA :\n\n👁️ = Support d'images | 📄 = Support PDF | 🧠 = Raisonnement | 🌐 = Recherche web\n❌ = Indisponible (jetons insuffisants/abonnement requis)",
		ModelUpdateSuccess:     "✅ Modèle mis à jour avec succès !",
		ModelImageNotSupported: "❌ Le modèle sélectionné ne prend pas en charge les images. Veuillez choisir un modèle différent ou envoyer un message texte.",
		ModelPDFNotSupported:   "❌ Le modèle sélectionné ne prend pas en charge les fichiers PDF. Veuillez choisir un modèle différent ou envoyer un message texte.",

		// Model Names
		ModelGeminiFlash:   "🚀 Gemini 2.5 Flash (Rapide et efficace pour réponses vives)",
		ModelGPT4o:         "🧠 GPT-4o (Le plus capable pour tâches complexes)",
		ModelGPT4oMini:     "⚡ GPT-4o Mini (Vitesse et performance équilibrées)",
		ModelClaude4Sonnet: "🎭 Claude 4 Sonnet (Écriture créative et analyse)",
		ModelGeminiPro:     "🌸 Gemini 2.5 Pro (Longs documents et contexte)",
		ModelO3Mini:        "🤖 OpenAI o3-mini (Modèle de raisonnement avancé)",
		ModelDeepSeekChat:  "🔬 DeepSeek Chat v3 (Recherche et programmation)",
		ModelDeepSeekR1:    "🔍 DeepSeek R1 (Raisonnement profond et logique)",

		// Queue
		QueueMessageQueued: "⏳ Votre message a été mis en file d'attente (position : %d). Je le traiterai après avoir terminé la réponse actuelle.",

		// Web Search
		ButtonWebSearchOn:             "🌐 Recherche web : ACTIVÉE",
		ButtonWebSearchOff:            "🌐 Recherche web : DÉSACTIVÉE",
		WebSearchSubscriptionRequired: "🔐 Web search requires an active subscription. Please subscribe to use this feature.",
		WebSearchEnabledNotification:  "🌐 Web search enabled! Each search query will cost 1 additional premium token.",

		// Profile
		ProfileTitle:              "👤 Votre Profil",
		ProfileTokenBalance:       "💰 Solde de Jetons :",
		ProfilePremiumTokens:      "🟡 Premium : %d jetons",
		ProfileRegularTokens:      "🔵 Régulier : %d jetons",
		ProfileInsufficientTokens: "❌ Jetons insuffisants. Vous avez besoin de %d jetons %s pour utiliser ce modèle.",

		// Subscription
		SubscriptionTitle:        "💳 Abonnement",
		SubscriptionMonthlyOffer: "🌟 Abonnement Premium Mensuel\n\n✨ Obtenez 1500 jetons réguliers + 100 jetons premium chaque mois !\n\nPrix : ⭐ 600 Étoiles Telegram par mois",
		SubscriptionBuyButton:    "💰 S'abonner pour ⭐ 600 Étoiles",
		SubscriptionSuccess:      "🎉 Abonnement activé ! Vous avez reçu 1500 jetons réguliers et 100 jetons premium.",
		SubscriptionFailed:       "❌ L'abonnement a échoué. Veuillez réessayer.",
		SubscriptionActiveInfo:   "✅ Vous avez un abonnement actif ! %d jours restants.",
		SubscriptionExpired:      "❌ Votre abonnement a expiré. Abonnez-vous à nouveau pour continuer à recevoir des jetons.",

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
		ButtonProfile:           "👤 Profil",
		ButtonSubscription:      "💳 Abonnement",
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
		ModelSelectTitle:       "🤖 KI-Modell auswählen:\n\n👁️ = Bildunterstützung | 📄 = PDF-Unterstützung | 🧠 = Denkvermögen | 🌐 = Web-Suche\n❌ = Nicht verfügbar (unzureichende Token/Abonnement erforderlich)",
		ModelUpdateSuccess:     "✅ Modell erfolgreich aktualisiert!",
		ModelImageNotSupported: "❌ Das ausgewählte Modell unterstützt keine Bilder. Bitte wählen Sie ein anderes Modell oder senden Sie eine Textnachricht.",
		ModelPDFNotSupported:   "❌ Das ausgewählte Modell unterstützt keine PDF-Dateien. Bitte wählen Sie ein anderes Modell oder senden Sie eine Textnachricht.",

		// Model Names
		ModelGeminiFlash:   "🚀 Gemini 2.5 Flash (Schnell und effizient für rasche Antworten)",
		ModelGPT4o:         "🧠 GPT-4o (Am fähigsten für komplexe Aufgaben)",
		ModelGPT4oMini:     "⚡ GPT-4o Mini (Ausgewogene Geschwindigkeit und Leistung)",
		ModelClaude4Sonnet: "🎭 Claude 4 Sonnet (Kreatives Schreiben und Analyse)",
		ModelGeminiPro:     "🌸 Gemini 2.5 Pro (Lange Dokumente und Kontext)",
		ModelO3Mini:        "🤖 OpenAI o3-mini (Fortgeschrittenes Denkmodell)",
		ModelDeepSeekChat:  "🔬 DeepSeek Chat v3 (Forschung und Programmierung)",
		ModelDeepSeekR1:    "🔍 DeepSeek R1 (Tiefes Denken und Logik)",

		// Queue
		QueueMessageQueued: "⏳ Ihre Nachricht wurde in die Warteschlange eingereiht (Position: %d). Ich werde sie nach Beendigung der aktuellen Antwort bearbeiten.",

		// Web Search
		ButtonWebSearchOn:             "🌐 Web-Suche: EIN",
		ButtonWebSearchOff:            "🌐 Web-Suche: AUS",
		WebSearchSubscriptionRequired: "🔐 Web search requires an active subscription. Please subscribe to use this feature.",
		WebSearchEnabledNotification:  "🌐 Web search enabled! Each search query will cost 1 additional premium token.",

		// Profile
		ProfileTitle:              "👤 Ihr Profil",
		ProfileTokenBalance:       "💰 Token-Guthaben:",
		ProfilePremiumTokens:      "🟡 Premium: %d Token",
		ProfileRegularTokens:      "🔵 Regulär: %d Token",
		ProfileInsufficientTokens: "❌ Unzureichende Token. Sie benötigen %d %s Token, um dieses Modell zu verwenden.",

		// Subscription
		SubscriptionTitle:        "💳 Abonnement",
		SubscriptionMonthlyOffer: "🌟 Monatliches Premium-Abonnement\n\n✨ Erhalten Sie jeden Monat 1500 reguläre Token + 100 Premium-Token!\n\nPreis: ⭐ 600 Telegram-Sterne pro Monat",
		SubscriptionBuyButton:    "💰 Abonnieren für ⭐ 600 Sterne",
		SubscriptionSuccess:      "🎉 Abonnement aktiviert! Sie haben 1500 reguläre Token und 100 Premium-Token erhalten.",
		SubscriptionFailed:       "❌ Abonnement fehlgeschlagen. Bitte versuchen Sie es erneut.",
		SubscriptionActiveInfo:   "✅ Sie haben ein aktives Abonnement! %d Tage verbleibend.",
		SubscriptionExpired:      "❌ Ihr Abonnement ist abgelaufen. Abonnieren Sie erneut, um weiterhin Token zu erhalten.",

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
		ButtonProfile:           "👤 Profilo",
		ButtonSubscription:      "💳 Abbonamento",
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
		ModelSelectTitle:       "🤖 Seleziona modello IA:\n\n👁️ = Supporto immagini | 📄 = Supporto PDF | 🧠 = Ragionamento | 🌐 = Ricerca web\n❌ = Non disponibile (token insufficienti/abbonamento richiesto)",
		ModelUpdateSuccess:     "✅ Modello aggiornato con successo!",
		ModelImageNotSupported: "❌ Il modello selezionato non supporta le immagini. Per favore scegli un modello diverso o invia un messaggio di testo.",
		ModelPDFNotSupported:   "❌ Il modello selezionato non supporta i file PDF. Per favore scegli un modello diverso o invia un messaggio di testo.",

		// Model Names
		ModelGeminiFlash:   "🚀 Gemini 2.5 Flash (Veloce ed efficiente per risposte rapide)",
		ModelGPT4o:         "🧠 GPT-4o (Più capace per compiti complessi)",
		ModelGPT4oMini:     "⚡ GPT-4o Mini (Velocità e prestazioni bilanciate)",
		ModelClaude4Sonnet: "🎭 Claude 4 Sonnet (Scrittura creativa e analisi)",
		ModelGeminiPro:     "🌸 Gemini 2.5 Pro (Documenti lunghi e contesto)",
		ModelO3Mini:        "🤖 OpenAI o3-mini (Modello di ragionamento avanzato)",
		ModelDeepSeekChat:  "🔬 DeepSeek Chat v3 (Ricerca e programmazione)",
		ModelDeepSeekR1:    "🔍 DeepSeek R1 (Ragionamento profondo e logica)",

		// Queue
		QueueMessageQueued: "⏳ Il tuo messaggio è stato messo in coda (posizione: %d). Lo elaborerò dopo aver terminato la risposta attuale.",

		// Web Search
		ButtonWebSearchOn:             "🌐 Ricerca web: ATTIVA",
		ButtonWebSearchOff:            "🌐 Ricerca web: DISATTIVA",
		WebSearchSubscriptionRequired: "🔐 Web search requires an active subscription. Please subscribe to use this feature.",
		WebSearchEnabledNotification:  "🌐 Web search enabled! Each search query will cost 1 additional premium token.",

		// Profile
		ProfileTitle:              "👤 Il Tuo Profilo",
		ProfileTokenBalance:       "💰 Saldo Token:",
		ProfilePremiumTokens:      "🟡 Premium: %d token",
		ProfileRegularTokens:      "🔵 Regolare: %d token",
		ProfileInsufficientTokens: "❌ Token insufficienti. Hai bisogno di %d token %s per usare questo modello.",

		// Subscription
		SubscriptionTitle:        "💳 Abbonamento",
		SubscriptionMonthlyOffer: "🌟 Abbonamento Premium Mensile\n\n✨ Ottieni 1500 token regolari + 100 token premium ogni mese!\n\nPrezzo: ⭐ 600 Stelle Telegram al mese",
		SubscriptionBuyButton:    "💰 Abbonati per ⭐ 600 Stelle",
		SubscriptionSuccess:      "🎉 Abbonamento attivato! Hai ricevuto 1500 token regolari e 100 token premium.",
		SubscriptionFailed:       "❌ Abbonamento fallito. Per favore riprova.",
		SubscriptionActiveInfo:   "✅ Hai un abbonamento attivo! %d giorni rimanenti.",
		SubscriptionExpired:      "❌ Il tuo abbonamento è scaduto. Abbonati di nuovo per continuare a ricevere token.",

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
		ButtonProfile:           "👤 个人资料",
		ButtonSubscription:      "💳 订阅",
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
		ModelSelectTitle:       "🤖 选择AI模型：\n\n👁️ = 图像支持 | 📄 = PDF支持 | 🧠 = 推理能力 | 🌐 = 网络搜索\n❌ = 不可用（代币不足/需要订阅）",
		ModelUpdateSuccess:     "✅ 模型更新成功！",
		ModelImageNotSupported: "❌ 所选模型不支持图像输入。请选择其他模型或发送文本消息。",
		ModelPDFNotSupported:   "❌ 所选模型不支持PDF文件。请选择其他模型或发送文本消息。",

		// Model Names
		ModelGeminiFlash:   "🚀 Gemini 2.5 Flash（快速高效，适合快速响应）",
		ModelGPT4o:         "🧠 GPT-4o（最强能力，适合复杂任务）",
		ModelGPT4oMini:     "⚡ GPT-4o Mini（平衡速度与性能）",
		ModelClaude4Sonnet: "🎭 Claude 4 Sonnet（创意写作与分析）",
		ModelGeminiPro:     "🌸 Gemini 2.5 Pro（长文档与上下文）",
		ModelO3Mini:        "🤖 OpenAI o3-mini（高级推理模型）",
		ModelDeepSeekChat:  "🔬 DeepSeek Chat v3（研究与编程）",
		ModelDeepSeekR1:    "🔍 DeepSeek R1（深度推理与逻辑）",

		// Queue
		QueueMessageQueued: "⏳ 您的消息已排队（位置：%d）。我会在完成当前回复后处理它。",

		// Web Search
		ButtonWebSearchOn:             "🌐 网络搜索：开启",
		ButtonWebSearchOff:            "🌐 网络搜索：关闭",
		WebSearchSubscriptionRequired: "🔐 Web search requires an active subscription. Please subscribe to use this feature.",
		WebSearchEnabledNotification:  "🌐 Web search enabled! Each search query will cost 1 additional premium token.",

		// Profile
		ProfileTitle:              "👤 您的个人资料",
		ProfileTokenBalance:       "💰 代币余额：",
		ProfilePremiumTokens:      "🟡 高级：%d 代币",
		ProfileRegularTokens:      "🔵 普通：%d 代币",
		ProfileInsufficientTokens: "❌ 代币不足。您需要 %d 个 %s 代币来使用此模型。",

		// Subscription
		SubscriptionTitle:        "💳 订阅",
		SubscriptionMonthlyOffer: "🌟 月度高级订阅\n\n✨ 每月获得 1500 个普通代币 + 100 个高级代币！\n\n价格：⭐ 每月 600 个 Telegram 星星",
		SubscriptionBuyButton:    "💰 订阅 ⭐ 600 星星",
		SubscriptionSuccess:      "🎉 订阅已激活！您已收到 1500 个普通代币和 100 个高级代币。",
		SubscriptionFailed:       "❌ 订阅失败。请重试。",
		SubscriptionActiveInfo:   "✅ 您有有效订阅！剩余 %d 天。",
		SubscriptionExpired:      "❌ 您的订阅已过期。请重新订阅以继续接收代币。",

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
		ButtonProfile:           "👤 プロフィール",
		ButtonSubscription:      "💳 サブスクリプション",
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
		ModelSelectTitle:       "🤖 AIモデルを選択：\n\n👁️ = 画像サポート | 📄 = PDFサポート | 🧠 = 推論機能 | 🌐 = ウェブ検索\n❌ = 利用不可（トークン不足/サブスクリプション必要）",
		ModelUpdateSuccess:     "✅ モデルが正常に更新されました！",
		ModelImageNotSupported: "❌ 選択されたモデルは画像入力をサポートしていません。別のモデルを選択するか、テキストメッセージを送信してください。",
		ModelPDFNotSupported:   "❌ 選択されたモデルはPDFファイルをサポートしていません。別のモデルを選択するか、テキストメッセージを送信してください。",

		// Model Names
		ModelGeminiFlash:   "🚀 Gemini 2.5 Flash（高速で効率的、素早い応答に最適）",
		ModelGPT4o:         "🧠 GPT-4o（最高性能、複雑なタスクに最適）",
		ModelGPT4oMini:     "⚡ GPT-4o Mini（速度と性能のバランス）",
		ModelClaude4Sonnet: "🎭 Claude 4 Sonnet（創造的な文章作成と分析）",
		ModelGeminiPro:     "🌸 Gemini 2.5 Pro（長文書とコンテキスト）",
		ModelO3Mini:        "🤖 OpenAI o3-mini（高度な推論モデル）",
		ModelDeepSeekChat:  "🔬 DeepSeek Chat v3（研究とプログラミング）",
		ModelDeepSeekR1:    "🔍 DeepSeek R1（深い推論と論理）",

		// Queue
		QueueMessageQueued: "⏳ メッセージがキューに追加されました（位置：%d）。現在の応答を完了した後に処理します。",

		// Web Search
		ButtonWebSearchOn:             "🌐 ウェブ検索：オン",
		ButtonWebSearchOff:            "🌐 ウェブ検索：オフ",
		WebSearchSubscriptionRequired: "🔐 Web search requires an active subscription. Please subscribe to use this feature.",
		WebSearchEnabledNotification:  "🌐 Web search enabled! Each search query will cost 1 additional premium token.",

		// Profile
		ProfileTitle:              "👤 あなたのプロフィール",
		ProfileTokenBalance:       "💰 トークン残高：",
		ProfilePremiumTokens:      "🟡 プレミアム：%d トークン",
		ProfileRegularTokens:      "🔵 通常：%d トークン",
		ProfileInsufficientTokens: "❌ トークンが不足しています。このモデルを使用するには %d 個の %s トークンが必要です。",

		// Subscription
		SubscriptionTitle:        "💳 サブスクリプション",
		SubscriptionMonthlyOffer: "🌟 月額プレミアムサブスクリプション\n\n✨ 毎月1500個の通常トークン + 100個のプレミアムトークンを取得！\n\n料金：⭐ 月額600テレグラムスター",
		SubscriptionBuyButton:    "💰 ⭐ 600スターで購読",
		SubscriptionSuccess:      "🎉 サブスクリプションが有効になりました！1500個の通常トークンと100個のプレミアムトークンを受け取りました。",
		SubscriptionFailed:       "❌ サブスクリプションに失敗しました。もう一度お試しください。",
		SubscriptionActiveInfo:   "✅ 有効なサブスクリプションがあります！残り %d 日。",
		SubscriptionExpired:      "❌ サブスクリプションの期限が切れました。トークンを受け取り続けるには、再度購読してください。",

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
		ButtonProfile:           "👤 프로필",
		ButtonSubscription:      "💳 구독",
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
		ModelSelectTitle:       "🤖 AI 모델 선택:\n\n👁️ = 이미지 지원 | 📄 = PDF 지원 | 🧠 = 추론 | 🌐 = 웹 검색\n❌ = 사용 불가 (토큰 부족/구독 필요)",
		ModelUpdateSuccess:     "✅ 모델이 성공적으로 업데이트되었습니다!",
		ModelImageNotSupported: "❌ 선택한 모델은 이미지 입력을 지원하지 않습니다. 다른 모델을 선택하거나 텍스트 메시지를 보내주세요.",
		ModelPDFNotSupported:   "❌ 선택한 모델은 PDF 파일을 지원하지 않습니다. 다른 모델을 선택하거나 텍스트 메시지를 보내주세요.",

		// Queue
		QueueMessageQueued: "⏳ 메시지가 대기열에 추가되었습니다 (위치: %d). 현재 응답을 완료한 후 처리하겠습니다.",

		// Web Search
		ButtonWebSearchOn:             "🌐 웹 검색: 켜짐",
		ButtonWebSearchOff:            "🌐 웹 검색: 꺼짐",
		WebSearchSubscriptionRequired: "🔐 Web search requires an active subscription. Please subscribe to use this feature.",
		WebSearchEnabledNotification:  "🌐 Web search enabled! Each search query will cost 1 additional premium token.",

		// Profile
		ProfileTitle:              "👤 프로필",
		ProfileTokenBalance:       "💰 토큰 잔고:",
		ProfilePremiumTokens:      "🟡 프리미엄: %d 토큰",
		ProfileRegularTokens:      "🔵 일반: %d 토큰",
		ProfileInsufficientTokens: "❌ 토큰이 부족합니다. 이 모델을 사용하려면 %d개의 %s 토큰이 필요합니다.",

		// Subscription
		SubscriptionTitle:        "💳 구독",
		SubscriptionMonthlyOffer: "🌟 월간 프리미엄 구독\n\n✨ 매달 1500개의 일반 토큰 + 100개의 프리미엄 토큰을 받으세요!\n\n가격: ⭐ 월 600 텔레그램 스타",
		SubscriptionBuyButton:    "💰 ⭐ 600 스타로 구독",
		SubscriptionSuccess:      "🎉 구독이 활성화되었습니다! 1500개의 일반 토큰과 100개의 프리미엄 토큰을 받았습니다.",
		SubscriptionFailed:       "❌ 구독 실패. 다시 시도해주세요.",
		SubscriptionActiveInfo:   "✅ 활성 구독이 있습니다! %d일 남음.",
		SubscriptionExpired:      "❌ 구독이 만료되었습니다. 토큰을 계속 받으려면 다시 구독하세요.",

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

		// Model Names
		ModelGeminiFlash:   "🚀 Gemini 2.5 Flash (빠르고 효율적인 응답)",
		ModelGPT4o:         "🧠 GPT-4o (복잡한 작업에 가장 적합)",
		ModelGPT4oMini:     "⚡ GPT-4o Mini (속도와 성능의 균형)",
		ModelClaude4Sonnet: "🎭 Claude 4 Sonnet (창의적 글쓰기 및 분석)",
		ModelGeminiPro:     "🌸 Gemini 2.5 Pro (긴 문서 및 컨텍스트)",
		ModelO3Mini:        "🤖 OpenAI o3-mini (고급 추론 모델)",
		ModelDeepSeekChat:  "🔬 DeepSeek Chat v3 (연구 및 코딩)",
		ModelDeepSeekR1:    "🔍 DeepSeek R1 (심층 추론 및 논리)",
	},
	"pt": {
		// Buttons
		ButtonStartConversation: "💬️ Iniciar Conversa",
		ButtonModelSelect:       "🤖 Selecionar Modelo",
		ButtonBackToMenu:        "🔙 Voltar ao Menu",
		ButtonNewConversation:   "➕ Nova Conversa",
		ButtonSettings:          "⚙️ Configurações",
		ButtonLanguage:          "🌐 Idioma",
		ButtonProfile:           "👤 Perfil",
		ButtonSubscription:      "💳 Assinatura",
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
		ModelSelectTitle:       "🤖 Selecionar modelo de IA:\n\n👁️ = Suporte a imagens | 📄 = Suporte a PDF | 🧠 = Raciocínio | 🌐 = Busca na web\n❌ = Não disponível (tokens insuficientes/assinatura necessária)",
		ModelUpdateSuccess:     "✅ Modelo atualizado com sucesso!",
		ModelImageNotSupported: "❌ O modelo selecionado não suporta entradas de imagem. Por favor, escolha um modelo diferente ou envie uma mensagem de texto.",
		ModelPDFNotSupported:   "❌ O modelo selecionado não suporta arquivos PDF. Por favor, escolha um modelo diferente ou envie uma mensagem de texto.",

		// Queue
		QueueMessageQueued: "⏳ Sua mensagem foi colocada na fila (posição: %d). Vou processá-la após terminar a resposta atual.",

		// Web Search
		ButtonWebSearchOn:             "🌐 Pesquisa web: LIGADA",
		ButtonWebSearchOff:            "🌐 Pesquisa web: DESLIGADA",
		WebSearchSubscriptionRequired: "🔐 Web search requires an active subscription. Please subscribe to use this feature.",
		WebSearchEnabledNotification:  "🌐 Web search enabled! Each search query will cost 1 additional premium token.",

		// Profile
		ProfileTitle:              "👤 Seu Perfil",
		ProfileTokenBalance:       "💰 Saldo de Tokens:",
		ProfilePremiumTokens:      "🟡 Premium: %d tokens",
		ProfileRegularTokens:      "🔵 Regular: %d tokens",
		ProfileInsufficientTokens: "❌ Tokens insuficientes. Você precisa de %d tokens %s para usar este modelo.",

		// Subscription
		SubscriptionTitle:        "💳 Assinatura",
		SubscriptionMonthlyOffer: "🌟 Assinatura Premium Mensal\n\n✨ Receba 1500 tokens regulares + 100 tokens premium todos os meses!\n\nPreço: ⭐ 600 Estrelas do Telegram por mês",
		SubscriptionBuyButton:    "💰 Assinar por ⭐ 600 Estrelas",
		SubscriptionSuccess:      "🎉 Assinatura ativada! Você recebeu 1500 tokens regulares e 100 tokens premium.",
		SubscriptionFailed:       "❌ A assinatura falhou. Por favor, tente novamente.",
		SubscriptionActiveInfo:   "✅ Você tem uma assinatura ativa! %d dias restantes.",
		SubscriptionExpired:      "❌ Sua assinatura expirou. Assine novamente para continuar recebendo tokens.",

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

		// Model Names
		ModelGeminiFlash:   "🚀 Gemini 2.5 Flash (Rápido e eficiente para respostas ágeis)",
		ModelGPT4o:         "🧠 GPT-4o (Mais capaz para tarefas complexas)",
		ModelGPT4oMini:     "⚡ GPT-4o Mini (Velocidade e desempenho equilibrados)",
		ModelClaude4Sonnet: "🎭 Claude 4 Sonnet (Escrita criativa e análise)",
		ModelGeminiPro:     "🌸 Gemini 2.5 Pro (Documentos longos e contexto)",
		ModelO3Mini:        "🤖 OpenAI o3-mini (Modelo de raciocínio avançado)",
		ModelDeepSeekChat:  "🔬 DeepSeek Chat v3 (Pesquisa e programação)",
		ModelDeepSeekR1:    "🔍 DeepSeek R1 (Raciocínio profundo e lógica)",
	},
	"hy": {
		// Buttons
		ButtonStartConversation: "💬️ Սկսել խոսակցություն",
		ButtonModelSelect:       "🤖 Ընտրել մոդելը",
		ButtonBackToMenu:        "🔙 Վերադառնալ մենյու",
		ButtonNewConversation:   "➕ Նոր խոսակցություն",
		ButtonSettings:          "⚙️ Կարգավորումներ",
		ButtonLanguage:          "🌐 Լեզու",
		ButtonProfile:           "👤 Պրոֆիլ",
		ButtonSubscription:      "💳 Բաժանորդագրություն",
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
		ModelSelectTitle:       "🤖 Ընտրել AI մոդելը:\n\n👁️ = Գծավոր աջակցություն | 📄 = PDF աջակցություն | 🧠 = Մտածողություն | 🌐 = Վեբ որոնում\n❌ = Հասանելի չէ (բավարար տոկեններ/բաժանորդագրություն պահանջ)",
		ModelUpdateSuccess:     "✅ Մոդելը հաջողությամբ թարմացվեց:",
		ModelImageNotSupported: "❌ Ընտրված մոդելը չի աջակցում պատկերների մուտքագրմանը: Խնդրում ենք ընտրել այլ մոդել կամ ուղարկել տեքստային հաղորդագրություն:",
		ModelPDFNotSupported:   "❌ Ընտրված մոդելը չի աջակցում PDF ֆայլերին: Խնդրում ենք ընտրել այլ մոդել կամ ուղարկել տեքստային հաղորդագրություն:",

		// Queue
		QueueMessageQueued: "⏳ Ձեր հաղորդագրությունը հերթի մեջ է (դիրքը՝ %d): Ես կմշակեմ այն ընթացիկ պատասխանն ավարտելուց հետո:",

		// Web Search
		ButtonWebSearchOn:             "🌐 Վեբ որոնում: ՄԻԱՑՎԱԾ",
		ButtonWebSearchOff:            "🌐 Վեբ որոնում: ԱՆՋԱՏՎԱԾ",
		WebSearchSubscriptionRequired: "🔐 Web search requires an active subscription. Please subscribe to use this feature.",
		WebSearchEnabledNotification:  "🌐 Web search enabled! Each search query will cost 1 additional premium token.",

		// Profile
		ProfileTitle:              "👤 Ձեր Պրոֆիլը",
		ProfileTokenBalance:       "💰 Տոկենների Մնացորդ:",
		ProfilePremiumTokens:      "🟡 Պրեմիում: %d տոկեն",
		ProfileRegularTokens:      "🔵 Սովորական: %d տոկեն",
		ProfileInsufficientTokens: "❌ Անբավարար տոկեններ: Այս մոդելն օգտագործելու համար ձեզ անհրաժեշտ է %d %s տոկեն:",

		// Subscription
		SubscriptionTitle:        "💳 Բաժանորդագրություն",
		SubscriptionMonthlyOffer: "🌟 Ամսական Պրեմիում Բաժանորդագրություն\n\n✨ Ստացեք 1500 սովորական տոկեն + 100 պրեմիում տոկեն ամեն ամիս!\n\nԳինը՝ ⭐ 600 Telegram աստղ ամսական",
		SubscriptionBuyButton:    "💰 Բաժանորդագրվել ⭐ 600 աստղով",
		SubscriptionSuccess:      "🎉 Բաժանորդագրությունն ակտիվացված է! Դուք ստացել եք 1500 սովորական տոկեն և 100 պրեմիում տոկեն:",
		SubscriptionFailed:       "❌ Բաժանորդագրությունը ձախողվեց: Խնդրում ենք փորձել կրկին:",
		SubscriptionActiveInfo:   "✅ Դուք ունեք ակտիվ բաժանորդագրություն! %d օր մնացել է:",
		SubscriptionExpired:      "❌ Ձեր բաժանորդագրությունը գործողության ժամկետն ավարտվել է: Տոկեններ ստանալու համար նորից բաժանորդագրվեք:",

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

		// Model Names
		ModelGeminiFlash:   "🚀 Gemini 2.5 Flash (Արագ և արդյունավետ պատասխանների համար)",
		ModelGPT4o:         "🧠 GPT-4o (Լավագույնը բարդ խնդիրների համար)",
		ModelGPT4oMini:     "⚡ GPT-4o Mini (Հավասարակշռված արագություն և արդյունավետություն)",
		ModelClaude4Sonnet: "🎭 Claude 4 Sonnet (Ստեղծագործական գրություն և վերլուծություն)",
		ModelGeminiPro:     "🌸 Gemini 2.5 Pro (Երկար փաստաթղթեր և համատեքստ)",
		ModelO3Mini:        "🤖 OpenAI o3-mini (Առաջադեմ տրամաբանության մոդել)",
		ModelDeepSeekChat:  "🔬 DeepSeek Chat v3 (Հետազոտություն և ծրագրավորում)",
		ModelDeepSeekR1:    "🔍 DeepSeek R1 (Խորը տրամաբանություն և տրամաբանություն)",
	},
	"uk": {
		// Buttons
		ButtonStartConversation: "💬️ Почати розмову",
		ButtonModelSelect:       "🤖 Вибрати модель",
		ButtonBackToMenu:        "🔙 Назад до меню",
		ButtonNewConversation:   "➕ Нова розмова",
		ButtonSettings:          "⚙️ Налаштування",
		ButtonLanguage:          "🌐 Мова",
		ButtonProfile:           "👤 Профіль",
		ButtonSubscription:      "💳 Підписка",
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
		ModelSelectTitle:       "🤖 Вибрати AI модель:\n\n👁️ = Підтримка зображень | 📄 = Підтримка PDF | 🧠 = Мислення | 🌐 = Пошук в інтернеті\n❌ = Недоступно (недостатньо токенів/потрібна підписка)",
		ModelUpdateSuccess:     "✅ Модель успішно оновлено!",
		ModelImageNotSupported: "❌ Обрана модель не підтримує зображення. Будь ласка, оберіть іншу модель або надішліть текстове повідомлення.",
		ModelPDFNotSupported:   "❌ Обрана модель не підтримує PDF файли. Будь ласка, оберіть іншу модель або надішліть текстове повідомлення.",

		// Queue
		QueueMessageQueued: "⏳ Ваше повідомлення поставлено в чергу (позиція: %d). Я оброблю його після завершення поточної відповіді.",

		// Web Search
		ButtonWebSearchOn:             "🌐 Веб-пошук: УВІМКНЕНО",
		ButtonWebSearchOff:            "🌐 Веб-пошук: ВИМКНЕНО",
		WebSearchSubscriptionRequired: "🔐 Web search requires an active subscription. Please subscribe to use this feature.",
		WebSearchEnabledNotification:  "🌐 Web search enabled! Each search query will cost 1 additional premium token.",

		// Profile
		ProfileTitle:              "👤 Ваш Профіль",
		ProfileTokenBalance:       "💰 Баланс Токенів:",
		ProfilePremiumTokens:      "🟡 Преміум: %d токенів",
		ProfileRegularTokens:      "🔵 Звичайні: %d токенів",
		ProfileInsufficientTokens: "❌ Недостатньо токенів. Вам потрібно %d %s токенів для використання цієї моделі.",

		// Subscription
		SubscriptionTitle:        "💳 Підписка",
		SubscriptionMonthlyOffer: "🌟 Щомісячна Преміум Підписка\n\n✨ Отримуйте 1500 звичайних токенів + 100 преміум токенів щомісяця!\n\nЦіна: ⭐ 600 Зірки Telegram на місяць",
		SubscriptionBuyButton:    "💰 Підписатися за ⭐ 600 Зірки",
		SubscriptionSuccess:      "🎉 Підписку активовано! Ви отримали 1500 звичайних токенів та 100 преміум токенів.",
		SubscriptionFailed:       "❌ Підписка не вдалася. Будь ласка, спробуйте знову.",
		SubscriptionActiveInfo:   "✅ У вас є активна підписка! Залишилося %d днів.",
		SubscriptionExpired:      "❌ Ваша підписка закінчилася. Підпишіться знову, щоб продовжити отримувати токени.",

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

		// Model Names
		ModelGeminiFlash:   "🚀 Gemini 2.5 Flash (Швидкий та ефективний для швидких відповідей)",
		ModelGPT4o:         "🧠 GPT-4o (Найкращий для складних завдань)",
		ModelGPT4oMini:     "⚡ GPT-4o Mini (Збалансована швидкість і продуктивність)",
		ModelClaude4Sonnet: "🎭 Claude 4 Sonnet (Творче письмо та аналіз)",
		ModelGeminiPro:     "🌸 Gemini 2.5 Pro (Довгі документи та контекст)",
		ModelO3Mini:        "🤖 OpenAI o3-mini (Модель розширеного міркування)",
		ModelDeepSeekChat:  "🔬 DeepSeek Chat v3 (Дослідження та програмування)",
		ModelDeepSeekR1:    "🔍 DeepSeek R1 (Глибоке міркування та логіка)",
	},
	"kk": {
		// Buttons
		ButtonStartConversation: "💬️ Сөйлесуді бастау",
		ButtonModelSelect:       "🤖 Модельді таңдау",
		ButtonBackToMenu:        "🔙 Мәзірге оралу",
		ButtonNewConversation:   "➕ Жаңа сөйлесу",
		ButtonSettings:          "⚙️ Баптаулар",
		ButtonLanguage:          "🌐 Тіл",
		ButtonProfile:           "👤 Профиль",
		ButtonSubscription:      "💳 Жазылым",
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
		ModelSelectTitle:       "🤖 AI моделін таңдау:\n\n👁️ = Кескін қолдауы | 📄 = PDF қолдауы | 🧠 = Ойлау | 🌐 = Веб іздеу\n❌ = Қолжетімсіз (токен жетіспейді/жазылым қажет)",
		ModelUpdateSuccess:     "✅ Модель сәтті жаңартылды!",
		ModelImageNotSupported: "❌ Таңдалған модель кескіндерді қолдамайды. Басқа модель таңдаңыз немесе мәтіндік хабарлама жіберіңіз.",
		ModelPDFNotSupported:   "❌ Таңдалған модель PDF файлдарын қолдамайды. Басқа модель таңдаңыз немесе мәтіндік хабарлама жіберіңіз.",

		// Queue
		QueueMessageQueued: "⏳ Сіздің хабарламаңыз кезекке қойылды (орын: %d). Мен оны ағымдағы жауапты аяқтағаннан кейін өңдеймін.",

		// Web Search
		ButtonWebSearchOn:             "🌐 Веб іздеу: ҚОСУЛЫ",
		ButtonWebSearchOff:            "🌐 Веб іздеу: ӨШІРУЛІ",
		WebSearchSubscriptionRequired: "🔐 Web search requires an active subscription. Please subscribe to use this feature.",
		WebSearchEnabledNotification:  "🌐 Web search enabled! Each search query will cost 1 additional premium token.",

		// Profile
		ProfileTitle:              "👤 Сіздің Профиліңіз",
		ProfileTokenBalance:       "💰 Токен Балансы:",
		ProfilePremiumTokens:      "🟡 Премиум: %d токен",
		ProfileRegularTokens:      "🔵 Кәдімгі: %d токен",
		ProfileInsufficientTokens: "❌ Токендер жеткіліксіз. Осы модельді пайдалану үшін сізге %d %s токен қажет.",

		// Subscription
		SubscriptionTitle:        "💳 Жазылым",
		SubscriptionMonthlyOffer: "🌟 Айлық Премиум Жазылым\n\n✨ Ай сайын 1500 қарапайым токен + 100 премиум токен алыңыз!\n\nБағасы: ⭐ Айына 600 Telegram жұлдызы",
		SubscriptionBuyButton:    "💰 ⭐ 600 жұлдызға жазылу",
		SubscriptionSuccess:      "🎉 Жазылым белсендірілді! Сіз 1500 қарапайым токен және 100 премиум токен алдыңыз.",
		SubscriptionFailed:       "❌ Жазылым сәтсіз болды. Қайта көріңіз.",
		SubscriptionActiveInfo:   "✅ Сізде белсенді жазылым бар! %d күн қалды.",
		SubscriptionExpired:      "❌ Сіздің жазылымыңыз аяқталды. Токендерді алуды жалғастыру үшін қайта жазылыңыз.",

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

		// Model Names
		ModelGeminiFlash:   "🚀 Gemini 2.5 Flash (Жылдам жауаптар үшін тез және тиімді)",
		ModelGPT4o:         "🧠 GPT-4o (Күрделі тапсырмалар үшін ең жақсы)",
		ModelGPT4oMini:     "⚡ GPT-4o Mini (Теңдестірілген жылдамдық және өнімділік)",
		ModelClaude4Sonnet: "🎭 Claude 4 Sonnet (Шығармашылық жазу және талдау)",
		ModelGeminiPro:     "🌸 Gemini 2.5 Pro (Ұзақ құжаттар және мәнмәтін)",
		ModelO3Mini:        "🤖 OpenAI o3-mini (Жетілдірілген ойлау моделі)",
		ModelDeepSeekChat:  "🔬 DeepSeek Chat v3 (Зерттеу және бағдарламалау)",
		ModelDeepSeekR1:    "🔍 DeepSeek R1 (Терең ойлау және логика)",
	},
	"ky": {
		// Buttons
		ButtonStartConversation: "💬️ Маекти баштоо",
		ButtonModelSelect:       "🤖 Модель тандоо",
		ButtonBackToMenu:        "🔙 Менюга кайтуу",
		ButtonNewConversation:   "➕ Жаңы маек",
		ButtonSettings:          "⚙️ Жөндөөлөр",
		ButtonLanguage:          "🌐 Тил",
		ButtonProfile:           "👤 Профиль",
		ButtonSubscription:      "💳 Жазылуу",
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
		ModelSelectTitle:       "🤖 AI модель тандоо:\n\n👁️ = Сүрөт колдуу | 📄 = PDF колдуу | 🧠 = Ой жүгүртүү | 🌐 = Веб издөө\n❌ = Жеткиликсиз (токен жетишпейт/жазылуу керек)",
		ModelUpdateSuccess:     "✅ Модель ийгиликтүү жаңыртылды!",
		ModelImageNotSupported: "❌ Тандалган модель сүрөттөрдү колдобойт. Башка модель тандаңыз же текст билдирүү жөнөтүңүз.",
		ModelPDFNotSupported:   "❌ Тандалган модель PDF файлдарды колдобойт. Башка модель тандаңыз же текст билдирүү жөнөтүңүз.",

		// Queue
		QueueMessageQueued: "⏳ Сиздин билдирүүңүз кезекке коюлду (орун: %d). Мен аны учурдагы жоопту бүткөндөн кийин иштетем.",

		// Web Search
		ButtonWebSearchOn:             "🌐 Веб издөө: КҮЙГҮЗҮЛГӨН",
		ButtonWebSearchOff:            "🌐 Веб издөө: ӨЧҮРҮЛГӨН",
		WebSearchSubscriptionRequired: "🔐 Web search requires an active subscription. Please subscribe to use this feature.",
		WebSearchEnabledNotification:  "🌐 Web search enabled! Each search query will cost 1 additional premium token.",

		// Profile
		ProfileTitle:              "👤 Сиздин Профилиңиз",
		ProfileTokenBalance:       "💰 Токен Баланс:",
		ProfilePremiumTokens:      "🟡 Премиум: %d токен",
		ProfileRegularTokens:      "🔵 Кадимки: %d токен",
		ProfileInsufficientTokens: "❌ Токендер жетишсиз. Бул модель колдонуу үчүн сизге %d %s токен керек.",

		// Subscription
		SubscriptionTitle:        "💳 Жазылуу",
		SubscriptionMonthlyOffer: "🌟 Айлык Премиум Жазылуу\\n\\n✨ Ар айда 1500 кадимки токен + 100 премиум токен алыңыз!\\n\\nБааси: ⭐ Айына 600 Telegram жылдызы",
		SubscriptionBuyButton:    "💰 ⭐ 600 жылдызга жазылуу",
		SubscriptionSuccess:      "🎉 Жазылуу иштетилди! Сиз 1500 кадимки токен жана 100 премиум токен алдыңыз.",
		SubscriptionFailed:       "❌ Жазылуу ийгиликсиз болду. Кайра аракет кылыңыз.",
		SubscriptionActiveInfo:   "✅ Сизде активдүү жазылуу бар! %d күн калды.",
		SubscriptionExpired:      "❌ Сиздин жазылууңуз бүттү. Токендерди алууну улантуу үчүн кайрадан жазылыңыз.",

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

		// Model Names
		ModelGeminiFlash:   "🚀 Gemini 2.5 Flash (Тез жана натыйжалуу жооптор үчүн)",
		ModelGPT4o:         "🧠 GPT-4o (Татаал тапшырмалар үчүн эң мыкты)",
		ModelGPT4oMini:     "⚡ GPT-4o Mini (Теңдештирилген ылдамдык жана аткаруу)",
		ModelClaude4Sonnet: "🎭 Claude 4 Sonnet (Чыгармачыл жазуу жана анализ)",
		ModelGeminiPro:     "🌸 Gemini 2.5 Pro (Узун документтер жана контекст)",
		ModelO3Mini:        "🤖 OpenAI o3-mini (Өркүндөтүлгөн ой жүгүртүү модели)",
		ModelDeepSeekChat:  "🔬 DeepSeek Chat v3 (Изилдөө жана программалоо)",
		ModelDeepSeekR1:    "🔍 DeepSeek R1 (Терең ой жүгүртүү жана логика)",
	},
	"ar": {
		// Buttons
		ButtonStartConversation: "💬️ بدء محادثة",
		ButtonModelSelect:       "🤖 اختيار نموذج",
		ButtonBackToMenu:        "🔙 العودة للقائمة",
		ButtonNewConversation:   "➕ محادثة جديدة",
		ButtonSettings:          "⚙️ الإعدادات",
		ButtonLanguage:          "🌐 اللغة",
		ButtonProfile:           "👤 الملف الشخصي",
		ButtonSubscription:      "💳 الاشتراك",
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
		ModelSelectTitle:       "🤖 اختيار نموذج الذكاء الاصطناعي:\n\n👁️ = دعم الصور | 📄 = دعم PDF | 🧠 = التفكير | 🌐 = بحث الويب\n❌ = غير متاح (رموز غير كافية/اشتراك مطلوب)",
		ModelUpdateSuccess:     "✅ تم تحديث النموذج بنجاح!",
		ModelImageNotSupported: "❌ النموذج المحدد لا يدعم مدخلات الصور. يرجى اختيار نموذج مختلف أو إرسال رسالة نصية.",
		ModelPDFNotSupported:   "❌ النموذج المحدد لا يدعم ملفات PDF. يرجى اختيار نموذج مختلف أو إرسال رسالة نصية.",

		// Queue
		QueueMessageQueued: "⏳ تم وضع رسالتك في طابور الانتظار (الموضع: %d). سأقوم بمعالجتها بعد إنهاء الرد الحالي.",

		// Web Search
		ButtonWebSearchOn:             "🌐 البحث على الويب: مُفعَّل",
		ButtonWebSearchOff:            "🌐 البحث على الويب: مُعطَّل",
		WebSearchSubscriptionRequired: "🔐 Web search requires an active subscription. Please subscribe to use this feature.",
		WebSearchEnabledNotification:  "🌐 Web search enabled! Each search query will cost 1 additional premium token.",

		// Profile
		ProfileTitle:              "👤 ملفك الشخصي",
		ProfileTokenBalance:       "💰 رصيد الرموز:",
		ProfilePremiumTokens:      "🟡 مميز: %d رمز",
		ProfileRegularTokens:      "🔵 عادي: %d رمز",
		ProfileInsufficientTokens: "❌ رموز غير كافية. تحتاج إلى %d رمز %s لاستخدام هذا النموذج.",

		// Subscription
		SubscriptionTitle:        "💳 الاشتراك",
		SubscriptionMonthlyOffer: "🌟 الاشتراك الشهري المميز\\n\\n✨ احصل على 1500 رمز عادي + 100 رمز مميز كل شهر!\\n\\nالسعر: ⭐ 600 نجمة تليجرام شهرياً",
		SubscriptionBuyButton:    "💰 اشترك مقابل ⭐ 600 نجمة",
		SubscriptionSuccess:      "🎉 تم تفعيل الاشتراك! لقد حصلت على 1500 رمز عادي و 100 رمز مميز.",
		SubscriptionFailed:       "❌ فشل الاشتراك. يرجى المحاولة مرة أخرى.",
		SubscriptionActiveInfo:   "✅ لديك اشتراك نشط! %d يوم متبقية.",
		SubscriptionExpired:      "❌ انتهت صلاحية اشتراكك. اشترك مرة أخرى لمواصلة تلقي الرموز.",

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

		// Model Names
		ModelGeminiFlash:   "🚀 Gemini 2.5 Flash (سريع وفعال للإجابات السريعة)",
		ModelGPT4o:         "🧠 GPT-4o (الأكثر قدرة للمهام المعقدة)",
		ModelGPT4oMini:     "⚡ GPT-4o Mini (سرعة وأداء متوازنان)",
		ModelClaude4Sonnet: "🎭 Claude 4 Sonnet (الكتابة الإبداعية والتحليل)",
		ModelGeminiPro:     "🌸 Gemini 2.5 Pro (مستندات طويلة وسياق)",
		ModelO3Mini:        "🤖 OpenAI o3-mini (نموذج تفكير متقدم)",
		ModelDeepSeekChat:  "🔬 DeepSeek Chat v3 (البحث والبرمجة)",
		ModelDeepSeekR1:    "🔍 DeepSeek R1 (التفكير العميق والمنطق)",
	},
	"hi": {
		// Buttons
		ButtonStartConversation: "💬️ बातचीत शुरू करें",
		ButtonModelSelect:       "🤖 मॉडल चुनें",
		ButtonBackToMenu:        "🔙 मेनू में वापस जाएं",
		ButtonNewConversation:   "➕ नई बातचीत",
		ButtonSettings:          "⚙️ सेटिंग्स",
		ButtonLanguage:          "🌐 भाषा",
		ButtonProfile:           "👤 प्रोफाइल",
		ButtonSubscription:      "💳 सब्सक्रिप्शन",
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
		ModelSelectTitle:       "🤖 AI मॉडल चुनें:\n\n👁️ = छवि समर्थन | 📄 = PDF समर्थन | 🧠 = तर्क | 🌐 = वेब खोज\n❌ = उपलब्ध नहीं (अपर्याप्त टोकन/सब्सक्रिप्शन आवश्यक)",
		ModelUpdateSuccess:     "✅ मॉडल सफलतापूर्वक अपडेट हो गया!",
		ModelImageNotSupported: "❌ चयनित मॉडल छवि इनपुट का समर्थन नहीं करता है। कृपया एक अलग मॉडल चुनें या टेक्स्ट संदेश भेजें।",
		ModelPDFNotSupported:   "❌ चयनित मॉडल PDF फाइलों का समर्थन नहीं करता है। कृपया एक अलग मॉडल चुनें या टेक्स्ट संदेश भेजें।",

		// Queue
		QueueMessageQueued: "⏳ आपका संदेश कतार में रखा गया है (स्थिति: %d)। मैं वर्तमान जवाब पूरा करने के बाد इसे प्रोसेस करूंगा।",

		// Web Search
		ButtonWebSearchOn:             "🌐 वेब खोज: चालू",
		ButtonWebSearchOff:            "🌐 वेब खोज: बंद",
		WebSearchSubscriptionRequired: "🔐 Web search requires an active subscription. Please subscribe to use this feature.",
		WebSearchEnabledNotification:  "🌐 Web search enabled! Each search query will cost 1 additional premium token.",

		// Profile
		ProfileTitle:              "👤 आपकी प्रोफाइल",
		ProfileTokenBalance:       "💰 टोकन बैलेंस:",
		ProfilePremiumTokens:      "🟡 प्रीमियम: %d टोकन",
		ProfileRegularTokens:      "🔵 नियमित: %d टोकन",
		ProfileInsufficientTokens: "❌ अपर्याप्त टोकन। इस मॉडल का उपयोग करने के लिए आपको %d %s टोकन की आवश्यकता है।",

		// Subscription
		SubscriptionTitle:        "💳 सब्सक्रिप्शन",
		SubscriptionMonthlyOffer: "🌟 मासिक प्रीमियम सब्सक्रिप्शन\\n\\n✨ हर महीने 1500 नियमित टोकन + 100 प्रीमियम टोकन पाएं!\\n\\nकीमत: ⭐ मासिक 600 टेलीग्राम स्टार",
		SubscriptionBuyButton:    "💰 ⭐ 600 स्टार के लिए सब्सक्राइब करें",
		SubscriptionSuccess:      "🎉 सब्सक्रिप्शन सक्रिय! आपको 1500 नियमित टोकन और 100 प्रीमियम टोकन मिले हैं।",
		SubscriptionFailed:       "❌ सब्सक्रिप्शन असफल। कृपया पुनः प्रयास करें।",
		SubscriptionActiveInfo:   "✅ आपके पास सक्रिय सब्सक्रिप्शन है! %d दिन शेष।",
		SubscriptionExpired:      "❌ आपकी सब्सक्रिप्शन समाप्त हो गई है। टोकन प्राप्त करना जारी रखने के लिए पुनः सब्सक्राइब करें।",

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

		// Model Names
		ModelGeminiFlash:   "🚀 Gemini 2.5 Flash (त्वरित उत्तरों के लिए तेज़ और कुशल)",
		ModelGPT4o:         "🧠 GPT-4o (जटिल कार्यों के लिए सबसे सक्षम)",
		ModelGPT4oMini:     "⚡ GPT-4o Mini (संतुलित गति और प्रदर्शन)",
		ModelClaude4Sonnet: "🎭 Claude 4 Sonnet (रचनात्मक लेखन और विश्लेषण)",
		ModelGeminiPro:     "🌸 Gemini 2.5 Pro (लंबे दस्तावेज़ और संदर्भ)",
		ModelO3Mini:        "🤖 OpenAI o3-mini (उन्नत तर्क मॉडल)",
		ModelDeepSeekChat:  "🔬 DeepSeek Chat v3 (अनुसंधान और कोडिंग)",
		ModelDeepSeekR1:    "🔍 DeepSeek R1 (गहन तर्क और तर्कशास्त्र)",
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
