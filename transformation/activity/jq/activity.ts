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

interface DefaultValues {
    [key: string]: any;
}

@WiContrib({})
@Injectable()
export class JQActivityContribution extends WiServiceHandlerContribution {

    // Constants for field names
    private readonly FIELD_NAMES = {
        ARGUMENTS: 'Arguments',
        ARGUMENT_NAMES: 'ArgumentNames',
        SCRIPT: 'Script'
    } as const;

    constructor(@Inject(Injector) injector, private http: Http) {
        super(injector, http);
    }

    private getDefaultValueForType(type: string): any {
        const defaults: DefaultValues = {
            String: "abc",
            Number: 0.1,
            Boolean: false,
            Object: {},
            Array: []
        };
        return defaults[type] ?? null;
    }

    value = (fieldName: string, context: IActivityContribution): Observable<any> | any => {
        if (fieldName === this.FIELD_NAMES.ARGUMENTS) {
            const argumentNames: IFieldDefinition = context.getField(this.FIELD_NAMES.ARGUMENT_NAMES);
            if (!!argumentNames.value && argumentNames.value.length > 0) {
                const jsonSchema: { [key: string]: any } = {};
                const data = JSON.parse(argumentNames.value);

                for (const item of data) {
                    jsonSchema[item.Name] = this.getDefaultValueForType(item.Type);
                }
                return JSON.stringify(jsonSchema);
            }
            return "{}";
        }
        return null;
    }

    validate = (fieldName: string, context: IActivityContribution): Observable<IValidationResult> | IValidationResult => {
        if (fieldName === this.FIELD_NAMES.ARGUMENTS) {
            const argumentNames: IFieldDefinition = context.getField(this.FIELD_NAMES.ARGUMENT_NAMES);
            const vresult: IValidationResult = ValidationResult.newValidationResult();
            vresult.setVisible(!!argumentNames.value && argumentNames.value.length > 0);
            return vresult;
        }
        if (fieldName === this.FIELD_NAMES.SCRIPT) {
            const script: IFieldDefinition = context.getField(this.FIELD_NAMES.SCRIPT);
            const vresult: IValidationResult = ValidationResult.newValidationResult();
            if (!!script.value && script.value !== "") {
                vresult.setValid(true);
            }
            vresult.setVisible(true);
            return vresult;
        }
    }
}
