package types

type EventResponse string

const (
	EventResponseSuccess               = EventResponse("response_success")
	EventResponseUnauthorized          = EventResponse("response_unauthorized")
	EventResponseBadRequest            = EventResponse("response_bad_request")
	EventResponseInternalError         = EventResponse("response_internal_error")
	EventResponseUnsupportedClientType = EventResponse("response_unsupported_client_type")
	EventResponseUnsupprtedEventType   = EventResponse("response_unsupported_event_type")
)

type EventList string

const (
	// Sample Events
	EventListSendMessage = EventList("send_message")

	// Business Logic Events
	EventListInitAuth                       = EventList("init_auth") // Initial Client On-Websocket-Connect Auth Event
	EventListSendPos                        = EventList("send_pos")  // Personnel sending updates their pos periodically
	EventListGetWsStats                     = EventList("get_ws_stats")
	EventListForwardOrderToDriver           = EventList("order_for_driver")             // API Forwards order to Drivers
	EventListUpdateZone                     = EventList("update_zone")                  // API Update Zones to Fleets
	EventListSendNotification               = EventList("send_notification")            // API Send Notification to Personnels and or Groups
	EventListSendAccountRegistrationRequest = EventList("account_registration_request") // API Send Account Registration Requests
	EventListGetAllPersonnelPosition        = EventList("get_all_personnel_pos")        // Command Center Get All Personnel Position
	EventListPanicButton                    = EventList("panic_button")                 // TBD
	EventListGetStakeholderMemberStats      = EventList("get_stakeholder_member_stats") // Get members online status for a stakeholder
	EventListPersonnelOnline                = EventList("personnel_online")             // Send updates to targets if a personnel is online
	EventListPersonnelOffline               = EventList("personnel_offline")            // Send updates to targets if a personnel is offline
)
