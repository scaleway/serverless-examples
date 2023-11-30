import pandas as pd
import numpy as np
from pydantic import BaseModel


class ClientProfile(BaseModel):
    """Represent features of client profile upon which the inference is applied"""

    age: int
    job: str
    marital: str
    education: str
    default: str
    housing: str
    loan: str
    contact: str
    month: str
    day_of_week: str
    duration: int
    campaign: int
    pdays: int
    previous: int
    poutcome: str
    emp_var_rate: float
    cons_price_idx: float
    cons_conf_idx: float
    euribor3m: float
    nr_employed: float


def clean_data(data: pd.DataFrame) -> pd.DataFrame:
    """Removes rows with missing value(s)"""

    data = data.dropna()
    return data


def transform_data(data: pd.DataFrame) -> pd.DataFrame:
    """
    This method handles the transformation of categorical variables of the dataset into 0/1 indicators.
    It also adds missing categorical variables that are by default false (0).
    """

    # # use the same category for basic education sub-categories
    data["education"] = np.where(
        data["education"] == "basic.9y", "Basic", data["education"]
    )
    data["education"] = np.where(
        data["education"] == "basic.6y", "Basic", data["education"]
    )
    data["education"] = np.where(
        data["education"] == "basic.4y", "Basic", data["education"]
    )

    # transform all categorical variables into 0/1 indicators and remove columns with string categories
    cat_vars = [
        "job",
        "marital",
        "education",
        "default",
        "housing",
        "loan",
        "contact",
        "month",
        "day_of_week",
        "poutcome",
    ]
    for var in cat_vars:
        cat_list = "var" + "_" + var
        cat_list = pd.get_dummies(data[var], prefix=var)
        data = data.join(cat_list)

    data_vars = data.columns.values.tolist()
    to_keep = [i for i in data_vars if i not in cat_vars]
    data = data[to_keep]

    # normalize column naming
    data.columns = data.columns.str.replace(".", "_")
    data.columns = data.columns.str.replace(" ", "_")

    # insert missing dummy categorical columns
    cat_vars_target = [
        "age",
        "duration",
        "campaign",
        "pdays",
        "previous",
        "emp_var_rate",
        "cons_price_idx",
        "cons_conf_idx",
        "euribor3m",
        "nr_employed",
        "job_admin_",
        "job_blue-collar",
        "job_entrepreneur",
        "job_housemaid",
        "job_management",
        "job_retired",
        "job_self-employed",
        "job_services",
        "job_student",
        "job_technician",
        "job_unemployed",
        "job_unknown",
        "marital_divorced",
        "marital_married",
        "marital_single",
        "marital_unknown",
        "education_Basic",
        "education_high_school",
        "education_illiterate",
        "education_professional_course",
        "education_university_degree",
        "education_unknown",
        "default_no",
        "default_unknown",
        "default_yes",
        "housing_no",
        "housing_unknown",
        "housing_yes",
        "loan_no",
        "loan_unknown",
        "loan_yes",
        "contact_cellular",
        "contact_telephone",
        "month_apr",
        "month_aug",
        "month_dec",
        "month_jul",
        "month_jun",
        "month_mar",
        "month_may",
        "month_nov",
        "month_oct",
        "month_sep",
        "day_of_week_fri",
        "day_of_week_mon",
        "day_of_week_thu",
        "day_of_week_tue",
        "day_of_week_wed",
        "poutcome_failure",
        "poutcome_nonexistent",
        "poutcome_success",
    ]

    for column_index, column_name in enumerate(cat_vars_target):
        if column_name not in data.columns:
            data.insert(column_index, column_name, False)

    return data
