// Wire types for POST /v0/requoter/onboard — mirror the Go dto4requoter shapes.

export interface IOnboardCar {
	name: string;
	/**
	 * Assetus asset id of a car the frontend already created/selected (via the
	 * Car step + REQUOTER_VEHICLE_SERVICE). When set, the onboard endpoint
	 * REFERENCES this asset instead of creating a new one from the detail fields
	 * below — the detail fields become advisory context. When empty, the backend
	 * creates the car server-side as before (backward compatible).
	 */
	assetID?: string;
	make?: string;
	model?: string;
	reg?: string;
	year?: number;
	prospective?: boolean;
	annualMileage?: number;
	classOfUse?: string;
}

export interface IOnboardDriver {
	firstName: string;
	lastName: string;
	isPrimary?: boolean;
}

export interface IOnboardInsurance {
	insurer?: string;
	policyNumber?: string;
	/** ISO YYYY-MM-DD renewal date. */
	renewalDate?: string;
}

export interface IOnboardRequest {
	spaceID: string;
	car: IOnboardCar;
	drivers: readonly IOnboardDriver[];
	insurance: IOnboardInsurance;
}

/** References to the records the onboard endpoint created (partial on failure). */
export interface IOnboardResponse {
	assetID?: string;
	contactIDs?: readonly string[];
	insuranceDocID?: string;
	happeningID?: string;
	profileID?: string;
	/** Set when a mid-sequence step failed; names the step. */
	failedStep?: string;
}
