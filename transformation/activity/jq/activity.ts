/*
 * Copyright © 2017. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
import { Observable } from "rxjs/Observable";
import { Injectable, Injector, Inject } from "@angular/core";
import { Http } from "@angular/http";
import {
    WiContrib,
    WiServiceHandlerContribution,
    IValidationResult,
    ValidationResult,
    IFieldDefinition,
    IActivityContribution,
    IConnectorContribution,
    WiContributionUtils
} from "wi-studio/app/contrib/wi-contrib";

@WiContrib({})
@Injectable()
export class JQActivityContribution extends WiServiceHandlerContribution {
    constructor(@Inject(Injector) injector, private http: Http) {
        super(injector, http);
    }

    value = (fieldName: string, context: IActivityContribution): Observable<any> | any => {
        if (fieldName === "Arguments") {
            let argumentNames: IFieldDefinition = context.getField("ArgumentNames");
            if (argumentNames.value) {
                // Read message attrbutes and construct JSON schema on the fly for the activity input
                var jsonSchema = {};
                // Convert string value into JSON object
                let data = JSON.parse(argumentNames.value);
                for (var i = 0; i < data.length; i++) {
                    if (data[i].Type === "String") {
                        jsonSchema[data[i].Name] = "abc";
                    } else if (data[i].Type === "Number") {
                        jsonSchema[data[i].Name] = 0.1;
                    }
                }
                return JSON.stringify(jsonSchema);
            }
            return "{}";
        }
        return null;
    }

    validate = (fieldName: string, context: IActivityContribution): Observable<IValidationResult> | IValidationResult => {
        if (fieldName === "Arguments") {
            let argumentNames: IFieldDefinition = context.getField("ArgumentNames");
            let vresult: IValidationResult = ValidationResult.newValidationResult();
            if (argumentNames.value) {
                vresult.setVisible(true);
            } else {
                vresult.setVisible(false);
            }
            return vresult;

        }
        if (fieldName === "Script") {
            let script: IFieldDefinition = context.getField("Script");
            let vresult: IValidationResult = ValidationResult.newValidationResult();
            if (script.value!=="") {
                vresult.setValid(true);
                vresult.setVisible(true);
            }
            return vresult;
        }
    }
}
